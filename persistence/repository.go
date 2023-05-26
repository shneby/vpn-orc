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
	DeletePeer(tenantId int, peerId string) error
}

func NewRepositoryService() RepositoryInterface {
	// real production scenario will have option to connect to a real database
	db, err := sql.Open("sqlite3", "resources/local.db")
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
		stmt.Close()
		return nil, err
	}

	row := stmt.QueryRow(sql.Named("tenantId", tenantId))

	tenant := &Tenant{}
	err = row.Scan(&tenant.Id, &tenant.Network)
	if err != nil {
		return nil, err
	}

	return tenant, stmt.Close()
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

	err := row.Scan(&peer.Id, &peer.VAddr, &peer.PublicKey, &peer.RAddr)
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
		rows.Close()
		return nil, err
	}

	var peers []Peer

	for rows.Next() {
		var peer Peer
		rows.Scan(&peer.Id, &peer.VAddr, &peer.PublicKey, &peer.RAddr)
		peers = append(peers, peer)
	}

	return peers, rows.Close()
}

func (r *RepositoryService) WritePeer(tenantId int, peer Peer) error {
	tableName := fmt.Sprintf("t%d_peers", tenantId)
	query := fmt.Sprintf("INSERT INTO '%s' (id, virtualAddress, publicKey, realAddress) VALUES (?, ?, ?, ?)", tableName)
	stmt, err := r.db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return err
	}

	_, err = stmt.Exec(peer.Id, peer.VAddr, peer.PublicKey, peer.RAddr)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryService) DeletePeer(tenantId int, peerId string) error {
	tableName := fmt.Sprintf("t%d_peers", tenantId)
	query := fmt.Sprintf("DELETE FROM '%s' WHERE id = :peerId", tableName)
	stmt, err := r.db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return err
	}

	_, err = stmt.Exec(sql.Named("peerId", peerId))
	if err != nil {
		return err
	}

	return nil
}
