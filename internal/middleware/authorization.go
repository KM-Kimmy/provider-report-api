package middleware

import (
	config "provider-report-api/configs"
	shared "provider-report-api/internal/modules/shared/dtos"
	"provider-report-api/pkg/utility"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// RoutePermission holds the required permissions for a route
type RoutePermission struct {
	Path       string
	Method     string
	Permission string
}

type Permission struct {
	MenuId   int      `json:"menuId"`
	MenuName string   `json:"menuName"`
	Actions  []string `json:"actions"`
}

// Example list of route permissions
var routePermissions = []RoutePermission{
	{Path: "/tpa-api/client-management/company", Method: "GET", Permission: "60:SELECT"},
	{Path: "/example", Method: "POST", Permission: "write:example"},
	// Add more route permissions as needed
}

// PermissionMiddleware checks if the user has the required permissions to access the endpoint
func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		claimsMap, ok := claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		if claimsMap["userRoleId"] == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"claimsMap": claimsMap, "claims": claims})
			c.Abort()
			return
		}

		userMasterId := claimsMap["sub"].(string)
		f, ok := claimsMap["userRoleId"].(float64)
		if !ok {
			fmt.Println("interface{} is not a float64")
			return
		}

		// Convert the float64 to a string
		userRoleId := strconv.FormatFloat(f, 'f', -1, 64)
		fmt.Println(userRoleId)
		var permissions *[]shared.UserActionAccessRights

		// Get the permission from Redis
		findRedisPermission, err := utility.NewRedisService().GetPermissions(userMasterId, userRoleId, "clientmanagement")
		if err != nil {
			permissions, err = FindUserAccessRights(userMasterId, userRoleId)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
				c.Abort()
				return
			}
			// save permission to redis
			err := utility.NewRedisService().SetPermissions(userMasterId, userRoleId, "clientmanagement", permissions)
			if err != nil {
				fmt.Println("error save permission to database")
			}
		} else {
			permissions = findRedisPermission
		}

		requiredPermission, exists := routeExists(c.FullPath(), c.Request.Method)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
			c.Abort()
			return
		}

		if !hasPermission(requiredPermission, permissions) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
			c.Abort()
			return
		}

		// If the user has permission, proceed to the next handler
		c.Next()
	}
}

// Function to check if permission exists
func hasPermission(requiredPermission string, permissions *[]shared.UserActionAccessRights) bool {
	parts := strings.Split(requiredPermission, ":")
	if len(parts) != 2 {
		return false
	}
	menuId := parts[0]
	action := parts[1]

	for _, perm := range *permissions {
		if fmt.Sprintf("%d", perm.MenuID) == menuId {
			for _, act := range perm.Actions {
				if act == action {
					return true
				}
			}
		}
	}
	return false
}

// Function to check if route exists and retrieve its required permission
func routeExists(path, method string) (string, bool) {
	for _, rp := range routePermissions {
		if rp.Path == path && rp.Method == method {
			return rp.Permission, true
		}
	}
	return "", false
}

func FindUserAccessRights(userMasterId, userRoleId string) (*[]shared.UserActionAccessRights, error) {
	db := config.GetDB() // Get the shared instance of the database

	query := `
		SELECT 
			m.MENU_ID,
			m.MENU_NAME,
			a.ACTION_ID,
			a.ACTION_CODE,
			a.ACTION_NAME
		FROM 
			tpacaredb.MAINTAIN.USER_ACCESS_RIGHTS uar
		LEFT JOIN 
			tpacaredb.MAINTAIN.USERS_USER_ACCESS_RIGHT_REL uuarr 
		ON 
			uar.USER_ACCESS_RIGHTS_ID = uuarr.USER_ACCESS_RIGHTS_ID
		LEFT JOIN 
			tpacaredb.MAINTAIN.USER_ACTION_ACCESS_RIGHTS uaar 
		ON 
			uar.USER_ACCESS_RIGHTS_ID = uaar.USER_ACCESS_RIGHTS_ID
		LEFT JOIN 
			tpacaredb.MAINTAIN.MENU_SCREEN_TAB_ACTION msta 
		ON 
			uaar.MENU_SCREEN_TAB_ACTION_ID = msta.MENU_SCREEN_TAB_ACTION_ID
		LEFT JOIN 
			tpacaredb.MAINTAIN.MENU m 
		ON 
			msta.MENU_ID = m.MENU_ID
		LEFT JOIN 
			tpacaredb.MAINTAIN.ACTION a 
		ON 
			msta.ACTION_ID = a.ACTION_ID
		WHERE 
			uuarr.USER_MASTER_ID = ? AND uar.USER_ROLE_ID = ? AND msta.STATUS_ID= 'Y' AND m.MENU_ID IN (?,?,?,?,?,?)`

	rows, err := db.Query(query, userMasterId, userRoleId, 60, 61, 5, 6, 7, 9)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menuMap = make(map[string]*shared.UserActionAccessRights)
	for rows.Next() {
		var menuID int
		var menuName string
		var action shared.UserAction
		err := rows.Scan(
			&menuID,
			&menuName,
			&action.ActionID,
			&action.ActionCode,
			&action.ActionName,
		)
		if err != nil {
			return nil, err
		}
		if menuAccess, exists := menuMap[menuName]; exists {
			actionExists := false
			if contains(menuAccess.Actions, *action.ActionCode) {
				actionExists = true
				break
			}
			if !actionExists {
				menuAccess.Actions = append(menuAccess.Actions, *action.ActionCode)
			}
		} else {
			menuAccess := &shared.UserActionAccessRights{
				MenuID:   menuID,
				MenuName: menuName,
				Actions:  []string{*action.ActionCode},
			}
			menuMap[menuName] = menuAccess
		}
	}
	finalData := convertMenuMapToSlice(menuMap)

	return &finalData, nil
}
func convertMenuMapToSlice(menuMap map[string]*shared.UserActionAccessRights) []shared.UserActionAccessRights {
	result := make([]shared.UserActionAccessRights, 0, len(menuMap))
	for _, value := range menuMap {
		result = append(result, *value)
	}
	return result
}

func contains(arr []string, target string) bool {
	seen := make(map[string]bool)
	for _, item := range arr {
		if seen[item] {
			continue
		}
		seen[item] = true
		if item == target {
			return true
		}
	}
	return false
}
