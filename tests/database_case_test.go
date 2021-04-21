package tests

import (
	"context"
	"myapp/graph/model"

	"github.com/stretchr/testify/require"
)

func (t *GormSuite) TestDatabaseCase() {
	//User Register
	var err error
	_, err = t.UserRegister(context.Background(), model.NewUser{
		Name:            "Dummy 1",
		Email:           "dummy1@gmail.com",
		Password:        "12345",
		ConfirmPassword: "12345",
	})
	require.NoError(t.T(), err)
	_, err = t.UserRegister(context.Background(), model.NewUser{
		Name:            "Dummy 2",
		Email:           "dummy2@gmail.com",
		Password:        "abcde",
		ConfirmPassword: "abcde",
	})
	require.NoError(t.T(), err)
	_, err = t.UserRegister(context.Background(), model.NewUser{
		Name:            "Dummy 3",
		Email:           "dummy3@gmail.com",
		Password:        "asdaqwezxdxzc",
		ConfirmPassword: "asdaqwezxdxzc",
	})
	require.NoError(t.T(), err)
	_, err = t.UserRegister(context.Background(), model.NewUser{
		Name:            "Dummy 4",
		Email:           "dummy4@gmail.com",
		Password:        "asdaqwezxdxzc",
		ConfirmPassword: "wrong password",
	})
	require.Error(t.T(), err)
	_, err = t.UserRegister(context.Background(), model.NewUser{
		Name:            "Dummy 4",
		Email:           "dummy1@gmail.com",
		Password:        "double email",
		ConfirmPassword: "double email",
	})
	require.Error(t.T(), err)

	//User Login
	_, err = t.UserLogin(context.Background(), "dummy1@gmail.com", "12345")
	require.NoError(t.T(), err)

	_, err = t.UserLogin(context.Background(), "dummy2@gmail.com", "abcde")
	require.NoError(t.T(), err)

	_, err = t.UserLogin(context.Background(), "dummy5@gmail.com", "unavailable email")
	require.Error(t.T(), err)

	_, err = t.UserLogin(context.Background(), "dummy3@gmail.com", "wrong password  ")
	require.Error(t.T(), err)

	//Team Create
	_, err = t.TeamCreate(context.Background(), 1, "Team Dummy User 1 A")
	require.NoError(t.T(), err)

	_, err = t.TeamCreate(context.Background(), 2, "Team Dummy User 2")
	require.NoError(t.T(), err)

	_, err = t.TeamCreate(context.Background(), 1, "Team Dummy User 1 B")
	require.NoError(t.T(), err)

	_, err = t.TeamCreate(context.Background(), 3, "Team Dummy User 3")
	require.NoError(t.T(), err)

	//Board Create
	_, err = t.BoardCreate(context.Background(), 1, model.NewBoard{
		Name:   "Board User 1 A",
		TeamID: 1,
	})
	require.NoError(t.T(), err)

	_, err = t.BoardCreate(context.Background(), 1, model.NewBoard{
		Name:   "Board User 1 B",
		TeamID: 3,
	})
	require.NoError(t.T(), err)

	_, err = t.BoardCreate(context.Background(), 3, model.NewBoard{
		Name:   "Board User 3 A",
		TeamID: 4,
	})
	require.NoError(t.T(), err)

	_, err = t.BoardCreate(context.Background(), 2, model.NewBoard{
		Name:   "Board Error",
		TeamID: 1,
	})
	require.Error(t.T(), err)

	//Team Assignment
	_, err = t.TeamAddMember(context.Background(), 2, model.NewTeamHasMember{
		TeamID: 1,
		UserID: 3,
	})
	require.Error(t.T(), err)

	_, err = t.TeamAddMember(context.Background(), 1, model.NewTeamHasMember{
		TeamID: 1,
		UserID: 2,
	})
	require.NoError(t.T(), err)

	_, err = t.TeamAddMember(context.Background(), 2, model.NewTeamHasMember{
		TeamID: 1,
		UserID: 3,
	})
	require.NoError(t.T(), err)

	_, err = t.TeamAddMember(context.Background(), 2, model.NewTeamHasMember{
		TeamID: 1,
		UserID: 3,
	})
	require.Error(t.T(), err)

	_, err = t.TeamRemoveMember(context.Background(), 1, model.NewTeamHasMember{
		TeamID: 1,
		UserID: 2,
	})
	require.NoError(t.T(), err)

	_, err = t.TeamRemoveMember(context.Background(), 1, model.NewTeamHasMember{
		TeamID: 1,
		UserID: 3,
	})
	require.NoError(t.T(), err)

	_, err = t.TeamAddMember(context.Background(), 2, model.NewTeamHasMember{
		TeamID: 1,
		UserID: 3,
	})
	require.Error(t.T(), err)

	//List Create
	_, err = t.ListCreateNext(context.Background(), 1, model.NewList{
		Name:    "List User 1 A",
		BoardID: 1,
	})
	require.NoError(t.T(), err)

	_, err = t.ListCreateNext(context.Background(), 1, model.NewList{
		Name:    "List User 1 B",
		BoardID: 1,
	})
	require.NoError(t.T(), err)

	_, err = t.ListCreateNext(context.Background(), 1, model.NewList{
		Name:    "List User 1 Board 2",
		BoardID: 2,
	})
	require.NoError(t.T(), err)

	_, err = t.ListCreateNext(context.Background(), 1, model.NewList{
		Name:    "List User 1 Board 2",
		BoardID: 3,
	})
	require.Error(t.T(), err)

	_, err = t.ListCreateNext(context.Background(), 2, model.NewList{
		Name:    "Failed To Create",
		BoardID: 1,
	})
	require.Error(t.T(), err)

	_, err = t.ListCreateNext(context.Background(), 3, model.NewList{
		Name:    "Failed To Create",
		BoardID: 1,
	})
	require.Error(t.T(), err)

	//List Item
	_, err = t.ListItemCreateNext(context.Background(), 1, model.NewListItem{
		Name:   "List Item User List 1 A",
		ListID: 1,
	})
	require.NoError(t.T(), err)

	_, err = t.ListItemCreateNext(context.Background(), 1, model.NewListItem{
		Name:   "List Item User List 1 B",
		ListID: 1,
	})
	require.NoError(t.T(), err)

	_, err = t.ListItemCreateNext(context.Background(), 1, model.NewListItem{
		Name:   "List Item User List 1 B",
		ListID: 2,
	})
	require.NoError(t.T(), err)

	_, err = t.ListItemCreateNext(context.Background(), 2, model.NewListItem{
		Name:   "List Item User List 1 B",
		ListID: 1,
	})
	require.Error(t.T(), err)

	_, err = t.ListItemCreateNext(context.Background(), 3, model.NewListItem{
		Name:   "List Item User List 1 B",
		ListID: 2,
	})
	require.Error(t.T(), err)

}
