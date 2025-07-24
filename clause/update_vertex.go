package clause

import (
	"errors"
	"fmt"
	"github.com/haysons/norm/resolver"
	"reflect"
	"sort"
	"strings"
)

type UpdateVertex struct {
	IsUpsert  bool
	VID       any
	TagUpdate any
	Opts      Options
}

const UpdateVertexName = "UPDATE_VERTEX"

func (uv UpdateVertex) Name() string {
	return UpdateVertexName
}

func (uv UpdateVertex) MergeIn(clause *Clause) {
	clause.Expression = uv
}

func (uv UpdateVertex) Build(nGQL Builder) error {
	if uv.IsUpsert {
		nGQL.WriteString("UPSERT VERTEX ON ")
	} else {
		nGQL.WriteString("UPDATE VERTEX ON ")
	}
	vidExpr, err := vertexIDExpr(uv.VID)
	if err != nil {
		return fmt.Errorf("norm: %w, build update vertex clause failed, %v", ErrInvalidClauseParams, err)
	}
	// name of the tag to be updated
	var tagName string
	tagNamer, ok := uv.TagUpdate.(resolver.VertexTagNamer)
	if ok {
		tagName = tagNamer.VertexTagName()
	}
	if uv.Opts.TagName != "" {
		tagName = uv.Opts.TagName
	}
	// list of properties to be updated
	propsName := make(map[string]bool, len(uv.Opts.PropNames))
	for _, propName := range uv.Opts.PropNames {
		propsName[propName] = true
	}
	propsUpdate, err := getPropsUpdateSet(uv.TagUpdate, propsName)
	if err != nil {
		return fmt.Errorf("norm: %w, build update vertex clause failed, %v", ErrInvalidClauseParams, err)
	}
	if vidExpr == "" {
		return fmt.Errorf("norm: %w, build update vertex clause failed, vid is empty", ErrInvalidClauseParams)
	}
	if tagName == "" {
		return fmt.Errorf("norm: %w, build update vertex clause failed, tag name is empty", ErrInvalidClauseParams)
	}
	if len(propsUpdate) == 0 {
		return fmt.Errorf("norm: %w, build update vertex clause failed, the values want to update empty", ErrInvalidClauseParams)
	}
	nGQL.WriteString(tagName)
	nGQL.WriteByte(' ')
	nGQL.WriteString(vidExpr)
	nGQL.WriteString(" SET ")
	for i, update := range propsUpdate {
		nGQL.WriteString(update[0])
		nGQL.WriteString(" = ")
		nGQL.WriteString(update[1])
		if i < len(propsUpdate)-1 {
			nGQL.WriteString(", ")
		}
	}
	return nil
}

func getPropsUpdateSet(propsUpdate any, needUpdate map[string]bool) ([][2]string, error) {
	propsUpdateSet := make([][2]string, 0)
	switch prop := propsUpdate.(type) {
	case map[string]any:
		for k, v := range prop {
			if len(needUpdate) > 0 && !needUpdate[k] {
				continue
			}
			propName := k
			var propValue string
			var err error
			switch expr := v.(type) {
			case Expr:
				exprBuilder := new(strings.Builder)
				err = expr.Build(exprBuilder)
				if err != nil {
					return nil, err
				}
				propValue = exprBuilder.String()
			case *Expr:
				exprBuilder := new(strings.Builder)
				err = expr.Build(exprBuilder)
				if err != nil {
					return nil, err
				}
				propValue = exprBuilder.String()
			default:
				propValue, err = resolver.FormatSimpleValue("", reflect.ValueOf(v))
				if err != nil {
					return nil, err
				}
			}
			propsUpdateSet = append(propsUpdateSet, [2]string{propName, propValue})
		}
		// sort by property name
		sort.Slice(propsUpdateSet, func(i, j int) bool {
			return propsUpdateSet[i][0] < propsUpdateSet[j][0]
		})
	default:
		propsValue := reflect.Indirect(reflect.ValueOf(propsUpdate))
		switch propsValue.Kind() {
		case reflect.Struct:
			propsType := propsValue.Type()
			for i := 0; i < propsType.NumField(); i++ {
				structField := propsType.Field(i)
				if structField.Anonymous || !structField.IsExported() {
					continue
				}
				propName := resolver.GetPropName(structField)
				sdkType := resolver.GetValueSdkType(structField)
				fieldValue := propsValue.Field(i)
				if len(needUpdate) > 0 && needUpdate[propName] {
					propValue, err := resolver.FormatSimpleValue(sdkType, fieldValue)
					if err != nil {
						return nil, err
					}
					propsUpdateSet = append(propsUpdateSet, [2]string{propName, propValue})
				} else if len(needUpdate) > 0 {
					continue
				} else {
					if fieldValue.IsZero() {
						continue
					}
					setting := resolver.ParseTagSetting(structField.Tag.Get(resolver.TagSettingKey))
					if setting[resolver.TagSettingIgnore] != "" || setting[resolver.TagSettingEdgeSrcID] != "" || setting[resolver.TagSettingEdgeDstID] != "" || setting[resolver.TagSettingEdgeRank] != "" || setting[resolver.TagSettingVertexID] != "" {
						continue
					}
					propValue, err := resolver.FormatSimpleValue(sdkType, fieldValue)
					if err != nil {
						return nil, err
					}
					propsUpdateSet = append(propsUpdateSet, [2]string{propName, propValue})
				}
			}
		case reflect.Map:
			propsType := propsValue.Type()
			if propsType.Key().Kind() != reflect.String {
				return nil, errors.New("update values must be map[string]any, struct or struct pointer")
			}
			updateMap := make(map[string]any)
			mapIter := propsValue.MapRange()
			for mapIter.Next() {
				k := mapIter.Key().String()
				v := mapIter.Value().Interface()
				updateMap[k] = v
			}
			return getPropsUpdateSet(updateMap, needUpdate)
		default:
			return nil, errors.New("update values must be map[string]any, struct or struct pointer")
		}
	}
	return propsUpdateSet, nil
}
