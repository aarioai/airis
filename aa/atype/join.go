package atype

import (
	"database/sql"
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/types"
	"golang.org/x/exp/constraints"
	"log"
	"reflect"
	"slices"
	"sort"
	"strings"
)

var (
	DefaultTagName = "name"
	JsonTagName    = "json"
	DbTagName      = "db"
)

type JoinType int

const (
	JoinSortedBit JoinType = 1 << 20
)

const (
	JoinKeys                  JoinType = 1 << iota // k
	JoinValues                                     // v
	JoinMySqlValues                                // "v"
	JoinKV                                         //kv
	JoinJSON                                       // "k":"v"
	JoinMySQL                                      // `t`.`k`="v"
	JoinMySqlFullLike                              // `t`.`k` LIKE "%v%"
	JoinMySqlStartWith                             // `t`.`k` LIKE "v%"
	JoinMySqlEndWith                               // `t`.`k` LIKE "%v"
	JoinMySqlLessThan                              // `t`.`k` < "v"
	JoinMySqlGreaterThan                           // `t`.`k` > "v"
	JoinMySqlGreaterOrEqualTo                      // `t`.`k` >= "v"
	JoinMySqlLessOrEqualTo                         // `t`.`k` <= "v"
	JoinURL                                        // k=v
	JoinSortedValues          = JoinSortedBit | JoinValues
	JoinSortedMySqlValues     = JoinSortedBit | JoinMySqlValues
	JoinSortedKV              = JoinSortedBit | JoinKV
	JoinSortedJSON            = JoinSortedBit | JoinJSON
	JoinSortedMySQL           = JoinSortedBit | JoinMySQL
	JoinSortedURL             = JoinSortedBit | JoinURL
)

var (
	mysqlJoinTypes = []JoinType{
		JoinMySqlValues, JoinMySQL, JoinMySqlFullLike, JoinMySqlStartWith, JoinMySqlEndWith,
		JoinMySqlLessThan, JoinMySqlGreaterThan, JoinMySqlGreaterOrEqualTo, JoinMySqlLessOrEqualTo,
		JoinSortedMySqlValues, JoinSortedMySQL,
	}
)

func toMySqlFieldName(k string) string {
	fields := strings.Split(k, ".")
	for i, field := range fields {
		fields[i] = "`" + field + "`"
	}
	return strings.Join(fields, ".")
}

func (t JoinType) TagName() string {
	if t == JoinJSON || t == JoinSortedJSON {
		return JsonTagName
	}
	if slices.Contains(mysqlJoinTypes, t) {
		return DbTagName
	}
	return DefaultTagName
}

func byJoinType(ty JoinType, k string, v any) string {
	var val string
	// @TODO separate sql from here
	if w, ok := v.(sql.NullBool); ok {
		val = String(w.Bool)
	} else if w, ok := v.(sql.NullFloat64); ok {
		val = String(w.Float64)
	} else if w, ok := v.(sql.NullInt64); ok {
		val = String(w.Int64)
	} else if w, ok := v.(sql.NullString); ok {
		val = w.String
	} else {
		val = String(v)
	}

	switch ty {
	case JoinKeys:
		return k
	case JoinSortedValues, JoinValues:
		return val
	case JoinSortedMySqlValues, JoinMySqlValues:
		return fmt.Sprintf(`"%s"`, val)
	case JoinSortedKV, JoinKV:
		return k + val
	case JoinSortedJSON, JoinJSON:
		return fmt.Sprintf(`"%s":"%s"`, k, val)
	case JoinSortedMySQL, JoinMySQL:
		return fmt.Sprintf(`%s="%s"`, toMySqlFieldName(k), val)
	case JoinMySqlFullLike:
		return fmt.Sprintf(`%s LIKE "%%%s%%"`, toMySqlFieldName(k), val)
	case JoinMySqlStartWith:
		return fmt.Sprintf(`%s LIKE "%s%%"`, toMySqlFieldName(k), val)
	case JoinMySqlEndWith:
		return fmt.Sprintf(`%s LIKE "%%%s"`, toMySqlFieldName(k), val)
	case JoinMySqlLessThan:
		return fmt.Sprintf(`%s<"%s"`, toMySqlFieldName(k), val)
	case JoinMySqlGreaterThan:
		return fmt.Sprintf(`%s>"%s"`, toMySqlFieldName(k), val)
	case JoinMySqlGreaterOrEqualTo:
		return fmt.Sprintf(`%s>="%s"`, toMySqlFieldName(k), val)
	case JoinMySqlLessOrEqualTo:
		return fmt.Sprintf(`%s<="%s"`, toMySqlFieldName(k), val)
	case JoinSortedURL, JoinURL:
		return fmt.Sprintf(`%s=%s`, k, val)
	}
	return ""
}

// JoinTagsByElements(stru, JoinUnsortedBit, " AND ", "json", "Name", "Age")
func JoinTagsByElements(u any, ty JoinType, sep string, tagname string, eles ...string) (ret string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	tags := make([]string, len(eles))

	t := reflect.TypeOf(u)
	for i, ele := range eles {
		for j := 0; j < t.NumField(); j++ {
			f := t.Field(j)
			if f.Name == ele {
				tags[i] = f.Tag.Get(tagname)
			}
		}
	}
	return JoinByTags(u, ty, sep, tagname, tags...)
}

// JoinByTags(stru, JoinUnsortedBit, " AND ", "json", "name", "age")
func JoinByTags(u any, ty JoinType, sep string, tagname string, tags ...string) (ret string) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("[error] atype.JoinByTags %s(%s %s): %s", sep, tagname, strings.Join(tags, ","), err)
		}
	}()

	if ty&JoinSortedBit > 0 {
		sort.Strings(tags)
	}

	t := reflect.TypeOf(u)
	var found bool
	for g := 0; g < len(tags); g++ {
		tag := tags[g]
		if tag == "" {
			continue
		}
		found = false
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			al := f.Tag.Get(tagname)
			if al == tag {
				found = true
				ret += sep + byJoinType(ty, tag, reflect.ValueOf(u).FieldByName(f.Name).Interface())
			}
		}
		if !found {
			ae.PanicF(`not found %s:"%s"`, tagname, tag)
		}
	}
	if len(ret) > len(sep) {
		ret = ret[len(sep):]
	}
	return
}

func JoinByNames(u any, ty JoinType, sep string, names ...string) string {
	return JoinByTags(u, ty, sep, ty.TagName(), names...)
}

func JoinNamesByElements(u any, ty JoinType, sep string, eles ...string) string {
	return JoinTagsByElements(u, ty, sep, ty.TagName(), eles...)
}

func JoinInt[T constraints.Signed](ids []T, sep string) string {
	if len(ids) == 0 {
		return ""
	}
	var s strings.Builder
	s.Grow(types.MaxInt64Len*len(ids) + ((len(ids) - 1) * len(sep)))
	for i, id := range ids {
		if i > 0 {
			s.WriteString(sep)
		}
		s.WriteString(types.FormatInt(id))
	}
	return s.String()
}
func JoinUint[T constraints.Unsigned](ids []T, sep string) string {
	if len(ids) == 0 {
		return ""
	}
	var s strings.Builder
	s.Grow(types.MaxUint64Len*len(ids) + ((len(ids) - 1) * len(sep)))
	for i, id := range ids {
		if i > 0 {
			s.WriteString(sep)
		}
		s.WriteString(types.FormatUint(id))
	}
	return s.String()
}
func JoinFloat[T constraints.Float](ids []T, sep string) string {
	if len(ids) == 0 {
		return ""
	}
	var s strings.Builder
	s.Grow(types.MaxUint64Len*len(ids) + ((len(ids) - 1) * len(sep)))
	for i, id := range ids {
		if i > 0 {
			s.WriteString(sep)
		}
		s.WriteString(types.FormatFloat(id))
	}
	return s.String()
}
