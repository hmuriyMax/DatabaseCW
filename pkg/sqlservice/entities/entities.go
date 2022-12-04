package entities

type TableNameStr struct {
	TableName string `db:"table_name"`
}

type Branch struct {
	BranchId  string `db:"branch_id" json:"branch_id"`
	ManagerId string `db:"manager_id" json:"manager_id"`
	Name      string `db:"name" json:"name"`
}

func (b Branch) TableName() string {
	return "branch"
}

type BranchList []Branch
