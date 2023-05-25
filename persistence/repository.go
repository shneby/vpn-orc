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
	ReadTenant(tenantId int) (*Tenant, error)
	ReadTenants() ([]Tenant, error)
	ReadPeer(tenantId int, peerId string) (*Peer, error)
	ReadPeers(tenantId int) ([]Peer, error)
	WritePeer(tenantId int, peer Peer) error
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

func (r *RepositoryService) ReadTenant(tenantId int) (*Tenant, error) {
	stmt, err := r.db.Prepare("SELECT * FROM tenants WHERE id = :tenantId")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	row := stmt.QueryRow(sql.Named("tenantId", tenantId))

	tenant := &Tenant{}
	err = row.Scan(&tenant.Id, &tenant.Network)
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

func (r *RepositoryService) ReadPeer(tenantId int, peerId string) (*Peer, error) {
	tableName := fmt.Sprintf("t%d_peers", tenantId)
	query := fmt.Sprintf("SELECT * FROM '%s' WHERE id = '%s'", tableName, peerId)
	row := r.db.QueryRow(query)
	peer := &Peer{}

	err := row.Scan(&peer.Id, &peer.Address, &peer.PublicKey)
	if err != nil {
		return nil, err
	}

	return peer, nil
}

func (r *RepositoryService) ReadPeers(tenantId int) ([]Peer, error) {
	tableName := fmt.Sprintf("t%d_peers", tenantId)
	query := fmt.Sprintf("SELECT * FROM '%s'", tableName)
	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var peers []Peer

	for rows.Next() {
		var peer Peer
		rows.Scan(&peer.Id, &peer.Address, &peer.PublicKey)
		peers = append(peers, peer)
	}

	return peers, nil
}

func (r *RepositoryService) WritePeer(tenantId int, peer Peer) error {
	tableName := fmt.Sprintf("t%d_peers", tenantId)
	query := fmt.Sprintf("INSERT INTO '%s' (id, address, publicKey) VALUES (?, ?, ?)", tableName)
	stmt, err := r.db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return err
	}

	_, err = stmt.Exec(peer.Id, peer.Address, peer.PublicKey)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryService) DeletePeer(tenantId int, peerId string) (*Peer, error) {
	return nil, nil
}
