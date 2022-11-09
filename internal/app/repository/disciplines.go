package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func GetDisciplineName(Db *sqlx.DB, disciplineID int) (string, error) {
	var disciplineName string

	err := Db.QueryRow("SELECT name from disciplines WHERE id =$1", disciplineID).Scan(&disciplineName)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Row with id unknown")
		} else {
			log.Println("Couldn't find the line with the order")
		}
		log.Error(err)
	}

	return disciplineName, err
}

//
//var ErrDisciplineNotFound = fmt.Errorf("discipline not found")
//
//func (r *Repository) SelectDiscipline(ctx context.Context, disciplineID int) (discipline *ds.Discipline, err error) {
//	// start ctx with timeout
//	disciplines := make([]ds.Discipline, 0)
//
//	ctx, cancel := context.WithTimeout(ctx, txTimeout)
//	defer cancel()
//
//	q := r.sq.Select().From("disciplines")
//	qSelect, argsSelect, err := q.Columns("id", "name").
//		Where(sq.Eq{"id": disciplineID}).
//		Limit(1).
//		ToSql()
//
//	err = r.db.SelectContext(ctx, &disciplines, qSelect, argsSelect...)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(disciplines) == 0 {
//		return nil, ErrDisciplineNotFound
//	}
//
//	return &disciplines[0], nil
//}
//
//func (r *Repository) SelectAllDisciplines(ctx context.Context) ([]ds.Discipline, error) {
//	var err error
//
//	disciplines := make([]ds.Discipline, 0)
//
//	// start ctx with timeout
//	ctx, cancel := context.WithTimeout(ctx, txTimeout)
//	defer cancel()
//
//	q := r.sq.Select().From("disciplines")
//	qSelect, argsSelect, err := q.Columns("id", "name").
//		OrderBy("id").
//		ToSql()
//
//	err = r.db.SelectContext(ctx, &disciplines, qSelect, argsSelect...)
//	if err != nil {
//		return nil, err
//	}
//
//	return disciplines, nil
//}
//
//func (r *Repository) SetPercentPerDisciplines(ctx context.Context, percent int, disciplines ...int) error {
//	var err error
//
//	// start ctx with timeout
//	ctx, cancel := context.WithTimeout(ctx, txTimeout)
//	defer cancel()
//
//	if disciplines == nil {
//		_, err = r.db.ExecContext(ctx, "UPDATE disciplines SET percent = $1 WHERE TRUE", percent)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func (r *Repository) AddNewDiscipline(ctx context.Context, d ds.Discipline) error {
//	var err error
//
//	// start ctx with timeout
//	ctx, cancel := context.WithTimeout(ctx, txTimeout)
//	defer cancel()
//
//	_, err = r.db.ExecContext(ctx, "INSERT INTO disciplines "+
//		"(name, percent) "+
//		"VALUES ($1, $2)", d.Name, d.Percent)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
