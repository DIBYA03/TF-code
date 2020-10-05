package services

import (
	"fmt"
	"reflect"
	"strings"
)

type SharedID string

func MaskLeft(s string, l int) string {
	rs := []rune(s)
	for i := 0; i < len(rs)-l; i++ {
		rs[i] = 'X'
	}

	return string(rs)
}

//SQLGenForUpdate Key Generator for db tags when updating
// it will generate a string with a pair of keys
// e.g `type User struct {
//Name  string `db:"name"`
//Email string `db:"email"`
//Optional *string `db:"optional"`
//Age int `db:"age"`
//}`
//Should ouput `name=:name,email=:email,optional=:optional,age=:age`
func SQLGenForUpdate(i interface{}) string {
	var query []string
	v := reflect.ValueOf(i)
	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Tag.Get("db")
		switch t := v.Field(i).Interface().(type) {
		default:
			if reflect.ValueOf(t).Kind() != reflect.Ptr {
				query = append(query, fmt.Sprintf("%s=:%s", key, key))
			} else if !reflect.ValueOf(t).IsNil() {
				query = append(query, fmt.Sprintf("%s=:%s", key, key))
			}

		}
	}
	return strings.Join(query, ",")
}

//SQLGenInsertKeys Key Generator for db tags when inserting
// it will generate a string with the key(field) names
// e.g `type User struct {
//Name  string `db:"name"`
//Email string `db:"email"`
//Optional *string `db:"optional"`
//Age int `db:"age"`
//}`
//Should ouput `name,email,optional,age`
func SQLGenInsertKeys(i interface{}) string {
	var query []string
	v := reflect.ValueOf(i)
	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Tag.Get("db")
		switch t := v.Field(i).Interface().(type) {
		default:
			if reflect.ValueOf(t).Kind() != reflect.Ptr {
				query = append(query, fmt.Sprintf("%s", key))
			} else if !reflect.ValueOf(t).IsNil() {
				query = append(query, fmt.Sprintf("%s", key))
			}

		}
	}
	return strings.Join(query, ",")
}

//SQLGenInsertValues Key Generator for db tags when inserting
// it will generate a string with a the value(field) names
// e.g `type User struct {
//Name  string `db:"name"`
//Email string `db:"email"`
//Optional *string `db:"optional"`
//Age int `db:"age"`
//}`
//Should ouput `:name,:email,:optional,:age`
func SQLGenInsertValues(i interface{}) string {
	var query []string
	v := reflect.ValueOf(i)
	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Tag.Get("db")
		switch t := v.Field(i).Interface().(type) {
		default:
			if reflect.ValueOf(t).Kind() != reflect.Ptr {
				query = append(query, fmt.Sprintf(":%s", key))
			} else if !reflect.ValueOf(t).IsNil() {
				query = append(query, fmt.Sprintf(":%s", key))
			}

		}
	}
	return strings.Join(query, ",")
}
