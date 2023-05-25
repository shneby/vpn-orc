package persistence

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type RepositoryService struct {
	db *sql.DB
}

type RepositoryInterface interface {
	ReadTenants() ([]Tenant, error)
	ReadTenant(tenantId string) (*Tenant, error)
}

func NewRepositoryService() RepositoryInterface {
	db, err := sql.Open("sqlite3", "local.db")
	if err != nil {
		log.Fatal("Failed to establish database connection")
	}

	return &RepositoryService{
		db: db,
	}
}

func (r *RepositoryService) ReadTenant(tenantId string) (*Tenant, error) {
	query := fmt.Sprintf("SELECT * FROM tenants WHERE id = '%s'", tenantId)
	row := r.db.QueryRow(query)
	tenant := &Tenant{}

	err := row.Scan(&tenant.Id, &tenant.Network)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (r *RepositoryService) ReadTenants() ([]Tenant, error) {
	query := fmt.Sprintf("SELECT * FROM tenants")
	rows, _ := r.db.Query(query)

	var tenants []Tenant

	for rows.Next() {
		var tenant Tenant
		rows.Scan(&tenant.Id, &tenant.Network)
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

func (r *RepositoryService) ReadPeer(tenantId string, peerId string) (*Peer, error) {
	return nil, nil
}

func (r *RepositoryService) WritePeer(tenantId string, peer Peer) {

}

func (r *RepositoryService) DeletePeer(tenantId string, peerId string) (*Peer, error) {
	return nil, nil
}