package graph

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/skeetcha/pf2ql/graph/model"
)

func GetData(rows *sql.Rows) (*model.Source, error) {
	var id string
	var name string
	var releaseDate string
	var productLine string
	var link string
	var errataVersion *float64
	var errataDate *string
	var isRemaster bool
	err := rows.Scan(&id, &name, &releaseDate, &productLine, &link, &errataVersion, &errataDate, &isRemaster)

	if err != nil {
		return nil, err
	}

	return &model.Source{ID: id, Name: name, ReleaseDate: releaseDate, ProductLine: model.ProductLine(productLine), Link: link, ErrataVersion: errataVersion, ErrataDate: errataDate, IsRemaster: isRemaster}, nil
}

type criterionData struct {
	Key string
	Modifier model.CriterionModifier
	Value any
}

func getModifierString(modifier model.CriterionModifier) string {
	if modifier == model.CriterionModifierEquals {
		return "="
	} else if modifier == model.CriterionModifierGreaterThan {
		return ">"
	} else if modifier == model.CriterionModifierLessThan {
		return "<"
	} else if modifier == model.CriterionModifierNotEquals {
		return "!="
	}

	return ""
}

func getCriteriaStr(criteria []criterionData) (string, error) {
	str := ""

	for _, v := range criteria {
		str += " "
		str += v.Key
		str += " "
		str += getModifierString(v.Modifier)
		str += " "

		if v.Key == "id" {
			val, ok := v.Value.(int)

			if ok {
				str += string(val)
			} else {
				return "", fmt.Errorf("%s is not an int", v.Key)
			}
		} else if v.Key == "releaseDate" || v.Key == "name" || v.Key == "link" || v.Key == "errataDate" {
			str += v.Value.(string)
		} else if v.Key == "productLine" {
			str += v.Value.(model.ProductLine).String()
		} else if v.Key == "errataVersion" {
			val, ok := v.Value.(float64)

			if ok {
				str += strconv.FormatFloat(val, 'f', -1, 64)
			} else {
				return "", fmt.Errorf("%s is not a float", v.Key)
			}
		} else if v.Key == "isRemaster" {
			val, ok := v.Value.(bool)

			if ok {
				if val {
					str += "TRUE"
				} else {
					str += "FALSE"
				}
			} else {
				return "", fmt.Errorf("%s is not a bool", v.Key)
			}
		}
	}

	return str, nil
}

func createCriteria(f *model.SourceFilter) []criterionData {
	criteria := []criterionData{}

	if f.ID != nil {
		criteria = append(criteria, criterionData {
			Key: "id",
			Modifier: f.ID.Modifier,
			Value: f.ID.Value,
		})
	}

	if f.Name != nil {
		criteria = append(criteria, criterionData {
			Key: "name",
			Modifier: f.Name.Modifier,
			Value: f.Name.Value,
		})
	}

	if f.ReleaseDate != nil {
		criteria = append(criteria, criterionData {
			Key: "releaseDate",
			Modifier: f.ReleaseDate.Modifier,
			Value: f.ReleaseDate.Value,
		})
	}

	if f.ProductLine != nil {
		criteria = append(criteria, criterionData {
			Key: "productLine",
			Modifier: f.ProductLine.Modifier,
			Value: f.ProductLine.Value,
		})
	}

	if f.Link != nil {
		criteria = append(criteria, criterionData {
			Key: "link",
			Modifier: f.Link.Modifier,
			Value: f.Link.Value,
		})
	}

	if f.ErrataVersion != nil {
		criteria = append(criteria, criterionData {
			Key: "errataVersion",
			Modifier: f.ErrataVersion.Modifier,
			Value: f.ErrataVersion.Value,
		})
	}

	if f.ErrataDate != nil {
		criteria = append(criteria, criterionData {
			Key: "errataDate",
			Modifier: f.ErrataDate.Modifier,
			Value: f.ErrataDate.Value,
		})
	}

	if f.IsRemaster != nil {
		criteria = append(criteria, criterionData {
			Key: "isRemaster",
			Modifier: f.IsRemaster.Modifier,
			Value: f.IsRemaster.Value,
		})
	}

	return criteria
}

func GetFilterString(f *model.SourceFilter) (string, error) {
	criteria := createCriteria(f)

	if f.And == nil && f.Or == nil && f.Not == nil {
		val, err := getCriteriaStr(criteria)

		if err != nil {
			return "", err
		}

		return val, nil
	} else {
		leftStr, err := getCriteriaStr(criteria)
		rightStr := ""

		if err != nil {
			return "", err
		}

		if f.And != nil {
			andStr, err := GetFilterString(f.And)

			if err != nil {
				return "", err
			}

			rightStr += " AND ("
			rightStr += andStr
			rightStr += ")"
		}

		if f.Or != nil {
			orStr, err := GetFilterString(f.Or)

			if err != nil {
				return "", err
			}

			rightStr += " OR ("
			rightStr += orStr
			rightStr += ")"
		}

		if f.Not != nil {
			notStr, err := GetFilterString(f.Not)

			if err != nil {
				return "", err
			}

			rightStr += " AND NOT ("
			rightStr += notStr
			rightStr += ")"
		}

		return leftStr + rightStr, nil
	}
}