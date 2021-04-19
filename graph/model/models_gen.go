// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AuthOps struct {
	Login    *JwtToken `json:"login"`
	Register *JwtToken `json:"register"`
}

type BoardOps struct {
	Create *Board `json:"create"`
}

type JwtToken struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type ListItemOps struct {
	Create *ListItem              `json:"create"`
	Move   map[string]interface{} `json:"move"`
}

type ListOps struct {
	Create *List   `json:"create"`
	Move   []*List `json:"move"`
}

type MoveList struct {
	ID           int  `json:"id"`
	MoveBeforeID *int `json:"move_before_id"`
	MoveAfterID  *int `json:"move_after_id"`
}

type MoveListItem struct {
	ID               int  `json:"id"`
	MoveBeforeID     *int `json:"move_before_id"`
	MoveAfterID      *int `json:"move_after_id"`
	MoveBeforeListID int  `json:"move_before_list_id"`
	MoveAfterListID  int  `json:"move_after_list_id"`
}

type NewBoard struct {
	Name   string `json:"name"`
	TeamID int    `json:"team_id"`
}

type NewList struct {
	Name    string `json:"name"`
	BoardID int    `json:"board_id"`
}

type NewListItem struct {
	Name   string `json:"name"`
	ListID int    `json:"list_id"`
}

type NewTeamHasMember struct {
	TeamID int `json:"team_id"`
	UserID int `json:"user_id"`
}

type NewUser struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type TeamHasMemberOps struct {
	Create *TeamHasMember `json:"create"`
}

type TeamOps struct {
	Create *Team `json:"create"`
}

type UserOps struct {
	EditName   string `json:"edit_name"`
	EditAvatar string `json:"edit_avatar"`
}
