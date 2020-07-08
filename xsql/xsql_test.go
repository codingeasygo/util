package xsql

import (
	"context"
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/jackc/pgx/v4/pgxpool"
)

var Pool func() *pgxpool.Pool

func init() {
	pool, err := pgxpool.Connect(context.Background(), "postgresql://cny:123@dev.loc:5432/dev")
	if err != nil {
		panic(err)
	}
	Pool = func() *pgxpool.Pool {
		return pool
	}
}

type timeObject struct {
	CreateTime Time
}

func TestTime(t *testing.T) {
	obj := &timeObject{
		CreateTime: Time(time.Now()),
	}
	bys, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}
	obj2 := &timeObject{}
	err = json.Unmarshal(bys, obj2)
	if err != nil {
		t.Error(err)
		return
	}
	if obj.CreateTime.Timestamp() != obj2.CreateTime.Timestamp() {
		t.Error("error")
		return
	}
	//
	t1 := TimeZero()
	bys, err = t1.MarshalJSON()
	if err != nil || string(bys) != "0" {
		t.Errorf("err:%v,bys:%v", err, string(bys))
		return
	}
	//
	//
	t2 := Time{}
	err = t2.UnmarshalJSON([]byte("null"))
	if err != nil {
		t.Error(err)
		return
	}
	bys, err = t2.MarshalJSON()
	if err != nil || string(bys) != "null" {
		t.Errorf("err:%v,bys:%v", err, string(bys))
		return
	}
	err = t2.Scan(time.Now())
	if err != nil {
		t.Error(err)
		return
	}
	//
	//
	TimeZero()
	TimeNow()
	TimeStartOfToday()
	TimeStartOfMonth()
	TimeUnix(0)
}

func TestIntArray(t *testing.T) {
	var ary IntArray
	err := ary.Scan(1)
	if err == nil {
		t.Error(err)
		return
	}
	err = ary.Scan("a")
	if err == nil {
		t.Error(err)
		return
	}
	// ary.Value()
	//
	v0, v1, v2 := int(3), int(2), int(1)
	ary = append(ary, v0)
	ary = append(ary, v1)
	ary = append(ary, v2)
	sort.Sort(ary)
	//
	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_int`)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = Pool().Exec(context.Background(), `create table xsql_test_int(tid int,iarry text)`)
	if err != nil {
		t.Error(err)
		return
	}
	{ //normal
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int values ($1,$2)`, 1, ary)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 IntArray
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int where tid=$1`, 1).Scan(&ary1)
		if err != nil || len(ary1) != 3 {
			t.Error(err)
			return
		}
		if !ary1.HavingOne(3) {
			t.Error("error")
			return
		}
		if ary1.HavingOne(4) {
			t.Error("error")
			return
		}
		if ary1.Join(",") != "1,2,3" {
			t.Error("error")
		}
		if ary1.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
		ary2 := append(ary1, 3).RemoveDuplicate()
		if ary2.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
	}
	{ //nil
		var arynil IntArray = nil
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int values ($1,$2)`, 2, arynil)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 IntArray
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int where tid=$1`, 2).Scan(&ary1)
		if err != nil || len(ary1) != 0 {
			t.Error(err)
			return
		}
	}
}

func TestIntPtrArray(t *testing.T) {
	var ary IntPtrArray
	err := ary.Scan(1)
	if err == nil {
		t.Error(err)
		return
	}
	err = ary.Scan("a")
	if err == nil {
		t.Error(err)
		return
	}
	// ary.Value()
	//
	v0, v1, v2 := int(3), int(2), int(1)
	ary = append(ary, &v0)
	ary = append(ary, &v1)
	ary = append(ary, &v2)
	sort.Sort(ary)
	//
	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_int`)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = Pool().Exec(context.Background(), `create table xsql_test_int(tid int,iarry text)`)
	if err != nil {
		t.Error(err)
		return
	}
	{ //normal
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int values ($1,$2)`, 1, ary)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 IntPtrArray
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int where tid=$1`, 1).Scan(&ary1)
		if err != nil || len(ary1) != 3 {
			t.Error(err)
			return
		}
		if !ary1.HavingOne(3) {
			t.Error("error")
			return
		}
		if ary1.HavingOne(4) {
			t.Error("error")
			return
		}
		if ary1.Join(",") != "1,2,3" {
			t.Error("error")
		}
		if ary1.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
		ary2 := append(ary1, converter.IntPtr(3)).RemoveDuplicate()
		if ary2.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
	}
	{ //nil
		var arynil IntPtrArray = nil
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int values ($1,$2)`, 2, arynil)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 IntPtrArray
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int where tid=$1`, 2).Scan(&ary1)
		if err != nil || len(ary1) != 0 {
			t.Error(err)
			return
		}
		{ //join
			var ary1 = IntPtrArray{}
			ary1 = append(ary1, &v0)
			ary1 = append(ary1, nil)
			ary1 = append(ary1, &v2)
			sort.Sort(ary1)
			if ary1.Join(",") != "nil,1,3" {
				t.Error(ary1.Join(","))
				return
			}
		}
	}
}

func TestInt64Array(t *testing.T) {
	var ary Int64Array
	err := ary.Scan(1)
	if err == nil {
		t.Error(err)
		return
	}
	err = ary.Scan("a")
	if err == nil {
		t.Error(err)
		return
	}
	// ary.Value()
	//
	v0, v1, v2 := int64(3), int64(2), int64(1)
	ary = append(ary, v0)
	ary = append(ary, v1)
	ary = append(ary, v2)
	sort.Sort(ary)
	//
	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_int64`)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = Pool().Exec(context.Background(), `create table xsql_test_int64(tid int,iarry text)`)
	if err != nil {
		t.Error(err)
		return
	}
	{ //normal
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 1, ary)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 Int64Array
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 1).Scan(&ary1)
		if err != nil || len(ary1) != 3 {
			t.Error(err)
			return
		}
		if !ary1.HavingOne(3) {
			t.Error("error")
			return
		}
		if ary1.HavingOne(4) {
			t.Error("error")
			return
		}
		if ary1.Join(",") != "1,2,3" {
			t.Error("error")
		}
		if ary1.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
		ary2 := append(ary1, 3).RemoveDuplicate()
		if ary2.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
	}
	{ //nil
		var arynil Int64Array = nil
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 2, arynil)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 Int64Array
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 2).Scan(&ary1)
		if err != nil || len(ary1) != 0 {
			t.Error(err)
			return
		}
	}
}

func TestInt64PtrArray(t *testing.T) {
	var ary Int64PtrArray
	err := ary.Scan(1)
	if err == nil {
		t.Error(err)
		return
	}
	err = ary.Scan("a")
	if err == nil {
		t.Error(err)
		return
	}
	// ary.Value()
	//
	v0, v1, v2 := int64(3), int64(2), int64(1)
	ary = append(ary, &v0)
	ary = append(ary, &v1)
	ary = append(ary, &v2)
	sort.Sort(ary)
	//
	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_int64`)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = Pool().Exec(context.Background(), `create table xsql_test_int64(tid int,iarry text)`)
	if err != nil {
		t.Error(err)
		return
	}
	{ //normal
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 1, ary)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 Int64PtrArray
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 1).Scan(&ary1)
		if err != nil || len(ary1) != 3 {
			t.Error(err)
			return
		}
		if !ary1.HavingOne(3) {
			t.Error("error")
			return
		}
		if ary1.HavingOne(4) {
			t.Error("error")
			return
		}
		if ary1.Join(",") != "1,2,3" {
			t.Error("error")
		}
		if ary1.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
		ary2 := append(ary1, converter.Int64Ptr(3)).RemoveDuplicate()
		if ary2.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
	}
	{ //nil
		var arynil Int64PtrArray = nil
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 2, arynil)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 Int64PtrArray
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 2).Scan(&ary1)
		if err != nil || len(ary1) != 0 {
			t.Error(err)
			return
		}
		{ //join
			var ary1 = Int64PtrArray{}
			ary1 = append(ary1, &v0)
			ary1 = append(ary1, nil)
			ary1 = append(ary1, &v2)
			sort.Sort(ary1)
			if ary1.Join(",") != "nil,1,3" {
				t.Error(ary1.Join(","))
				return
			}
		}
	}
}

func TestFloat64Array(t *testing.T) {
	var ary Float64Array
	err := ary.Scan(1)
	if err == nil {
		t.Error(err)
		return
	}
	err = ary.Scan("a")
	if err == nil {
		t.Error(err)
		return
	}
	// ary.Value()
	//
	v0, v1, v2 := float64(3), float64(2), float64(1)
	ary = append(ary, v0)
	ary = append(ary, v1)
	ary = append(ary, v2)
	sort.Sort(ary)
	//
	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_int64`)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = Pool().Exec(context.Background(), `create table xsql_test_int64(tid int,iarry text)`)
	if err != nil {
		t.Error(err)
		return
	}
	{ //normal
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 1, ary)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 Float64Array
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 1).Scan(&ary1)
		if err != nil || len(ary1) != 3 {
			t.Error(err)
			return
		}
		if !ary1.HavingOne(3) {
			t.Error("error")
			return
		}
		if ary1.HavingOne(4) {
			t.Error("error")
			return
		}
		if ary1.Join(",") != "1,2,3" {
			t.Error("error")
		}
		if ary1.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
		ary2 := append(ary1, 3).RemoveDuplicate()
		if ary2.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
	}
	{ //nil
		var arynil Float64Array = nil
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 2, arynil)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 Float64Array
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 2).Scan(&ary1)
		if err != nil || len(ary1) != 0 {
			t.Error(err)
			return
		}
	}
}

func TestFloat64PtrArray(t *testing.T) {
	var ary Float64PtrArray
	err := ary.Scan(1)
	if err == nil {
		t.Error(err)
		return
	}
	err = ary.Scan("a")
	if err == nil {
		t.Error(err)
		return
	}
	// ary.Value()
	//
	v0, v1, v2 := float64(3), float64(2), float64(1)
	ary = append(ary, &v0)
	ary = append(ary, &v1)
	ary = append(ary, &v2)
	sort.Sort(ary)
	//
	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_int64`)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = Pool().Exec(context.Background(), `create table xsql_test_int64(tid int,iarry text)`)
	if err != nil {
		t.Error(err)
		return
	}
	{ //normal
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 1, ary)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 Float64PtrArray
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 1).Scan(&ary1)
		if err != nil || len(ary1) != 3 {
			t.Error(err)
			return
		}
		if !ary1.HavingOne(3) {
			t.Error("error")
			return
		}
		if ary1.HavingOne(4) {
			t.Error("error")
			return
		}
		if ary1.Join(",") != "1,2,3" {
			t.Error(ary1.Join(","))
		}
		if ary1.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
		ary2 := append(ary1, converter.Float64Ptr(3)).RemoveDuplicate()
		if ary2.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
	}
	{ //nil
		var arynil Float64PtrArray = nil
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 2, arynil)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 Float64PtrArray
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 2).Scan(&ary1)
		if err != nil || len(ary1) != 0 {
			t.Error(err)
			return
		}
	}
	{ //join
		var ary1 = Float64PtrArray{}
		ary1 = append(ary1, &v0)
		ary1 = append(ary1, nil)
		ary1 = append(ary1, &v2)
		sort.Sort(ary1)
		if ary1.Join(",") != "nil,1,3" {
			t.Error(ary1.Join(","))
			return
		}
	}
}

func TestMap(t *testing.T) {
	var eval M
	err := eval.Scan(1)
	if err == nil {
		t.Error(err)
		return
	}
	//
	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_map`)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = Pool().Exec(context.Background(), `create table xsql_test_map(tid int,mval text)`)
	if err != nil {
		t.Error(err)
		return
	}
	{ //normal
		var mval = M{"a": 1}
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_map values ($1,$2)`, 1, mval)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var mval1 M
		err = Pool().QueryRow(context.Background(), `select mval from xsql_test_map where tid=$1`, 1).Scan(&mval1)
		if err != nil || len(mval1) != 1 {
			t.Error(err)
			return
		}
		mv1 := xmap.Wrap(xmap.M(mval1))
		if mv1.Int("a") != 1 {
			t.Error("error")
		}
	}
	{ //nil
		var mval M = nil
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_map values ($1,$2)`, 2, mval)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var mval1 M
		err = Pool().QueryRow(context.Background(), `select mval from xsql_test_map where tid=$1`, 2).Scan(&mval1)
		if err != nil || len(mval1) != 0 {
			t.Error(err)
			return
		}
		mval1.RawMap()
	}
}

func TestStringArray(t *testing.T) {
	var ary StringArray
	err := ary.Scan(1)
	if err == nil {
		t.Error(err)
		return
	}
	err = ary.Scan("a")
	if err == nil {
		t.Error(err)
		return
	}
	// ary.Value()
	//
	v0, v1, v2 := "3", "2", "1"
	ary = append(ary, v0)
	ary = append(ary, v1)
	ary = append(ary, v2)
	sort.Sort(ary)
	//
	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_string`)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = Pool().Exec(context.Background(), `create table xsql_test_string(tid int,sarry text)`)
	if err != nil {
		t.Error(err)
		return
	}
	{ //normal
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_string values ($1,$2)`, 1, ary)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 StringArray
		err = Pool().QueryRow(context.Background(), `select sarry from xsql_test_string where tid=$1`, 1).Scan(&ary1)
		if err != nil || len(ary1) != 3 {
			t.Error(err)
			return
		}
		if !ary1.HavingOne("3") {
			t.Error("error")
			return
		}
		if ary1.HavingOne("4") {
			t.Error("error")
			return
		}
		if ary1.Join(",") != "1,2,3" {
			t.Error("error")
		}
		if ary1.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
		ary2 := append(ary1, "3", "", " ").RemoveDuplicate(true, true)
		if ary2.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
		ary3 := append(ary1, "", " ").RemoveEmpty(true)
		if ary3.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
	}
	{ //nil
		var arynil StringArray = nil
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 2, arynil)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 StringArray
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 2).Scan(&ary1)
		if err != nil || len(ary1) != 0 {
			t.Error(err)
			return
		}
	}
}

func TestStringPtrArray(t *testing.T) {
	var ary StringPtrArray
	err := ary.Scan(1)
	if err == nil {
		t.Error(err)
		return
	}
	err = ary.Scan("a")
	if err == nil {
		t.Error(err)
		return
	}
	// ary.Value()
	//
	v0, v1, v2 := "3", "2", "1"
	ary = append(ary, &v0)
	ary = append(ary, &v1)
	ary = append(ary, &v2)
	sort.Sort(ary)
	//
	_, err = Pool().Exec(context.Background(), `drop table if exists xsql_test_string`)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = Pool().Exec(context.Background(), `create table xsql_test_string(tid int,sarry text)`)
	if err != nil {
		t.Error(err)
		return
	}
	{ //normal
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_string values ($1,$2)`, 1, ary)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 StringPtrArray
		err = Pool().QueryRow(context.Background(), `select sarry from xsql_test_string where tid=$1`, 1).Scan(&ary1)
		if err != nil || len(ary1) != 3 {
			t.Error(err)
			return
		}
		if !ary1.HavingOne("3") {
			t.Error("error")
			return
		}
		if ary1.HavingOne("4") {
			t.Error("error")
			return
		}
		if ary1.Join(",") != "1,2,3" {
			t.Error("error")
		}
		if ary1.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
		ary2 := append(ary1, nil, converter.StringPtr("3"), converter.StringPtr(""), converter.StringPtr(" ")).RemoveDuplicate(true, true)
		if ary2.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
		ary3 := append(ary1, nil, converter.StringPtr(""), converter.StringPtr(" ")).RemoveEmpty(true)
		if ary3.DbArray() != "{1,2,3}" {
			t.Error("error")
		}
	}
	{ //nil
		var arynil StringPtrArray = nil
		res, err := Pool().Exec(context.Background(), `insert into xsql_test_int64 values ($1,$2)`, 2, arynil)
		if err != nil || !res.Insert() {
			t.Error(err)
			return
		}
		var ary1 StringPtrArray
		err = Pool().QueryRow(context.Background(), `select iarry from xsql_test_int64 where tid=$1`, 2).Scan(&ary1)
		if err != nil || len(ary1) != 0 {
			t.Error(err)
			return
		}
	}
	{ //join
		var ary1 = StringPtrArray{}
		ary1 = append(ary1, &v0)
		ary1 = append(ary1, nil)
		ary1 = append(ary1, &v2)
		sort.Sort(ary1)
		if ary1.Join(",") != "nil,1,3" {
			t.Error(ary1.Join(","))
			return
		}
	}
}
