package shared

type UserAction struct {
	ActionID     *int    `json:"actionId"`
	ActionCode *string 	 `json:"actionCode"`
	ActionName   *string `json:"actionName"`
}

type UserActionAccessRights struct {
	MenuID   int      `json:"menuId"`
	MenuName string   `json:"menuName"`
	Actions  []string `json:"actions"`
}
