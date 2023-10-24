package repository

import (
	"context"
	"database/sql"
	"log"
	r "notification/domain/repository"
	"strconv"
	"strings"
)

type billingRepo struct {
	Conn *sql.DB
}

func NewRepository(conn *sql.DB) r.BillingRepository {
	return &billingRepo{
		Conn: conn,
	}
}

func (p *billingRepo)AddConsume(ctx context.Context,d []r.Consumo){
	query := "insert into consumo(profile_id,amount,type_entity,id_entity,message) values"
	vals := []interface{}{}
	filedNumbers := 5
	for i, val := range d {
		n := i * filedNumbers
		query += `(`
		for j := 0; j < filedNumbers; j++ {
			query += `$`+strconv.Itoa(n+j+1) + `,`
		}
		query = query[:len(query)-1] + `),`
		vals = 	append(vals,val.ProfileId,val.Amount,val.TypeEntity,val.IdEnitity,val.Message)
	}
	query = strings.TrimSuffix(query,",")
	_,err := p.Conn.ExecContext(ctx,query,vals...)
	if err != nil {
		log.Println(err)
	}
}