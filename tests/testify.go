package tests

import (
	"myapp/config"
	"myapp/graph/model"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

//getModels Get Model for Database Mock
func getModels() []interface{} {
	return []interface{}{&model.User{}, &model.Team{}, &model.Board{}, &model.List{}, &model.ListItem{}, &model.TeamHasMember{}}
}

//GormSuite Testing Suite
type GormSuite struct {
	suite.Suite
	db *gorm.DB
	tr *gorm.DB
}

func (t *GormSuite) SetupSuite() {
	t.db = config.ConnectGormTest()
	for _, val := range getModels() {
		t.db.AutoMigrate(val)
	}
}

func (t *GormSuite) TearDownSuite() {
	sqlDB, _ := t.db.DB()
	defer sqlDB.Close()
	for _, val := range getModels() {
		t.db.Migrator().DropTable(val)
	}
}

func (t *GormSuite) SetupTest() {
    t.tr  = t.db.Begin()
}

func (t *GormSuite) TearDownTest() {
    t.tr.Rollback()
}