package GABAHTMLParser

import (
	"strings"
)

type Element struct {
	Parent      *Element
	Child       []*Element
	Tag         string
	InnerHTML   string
	Attr        map[string]string
}

func (element *Element) Find(condition string) []*Element{
	results := []*Element{}
	conditions := strings.Split(condition,"&&")
	for _, child := range element.Child {
		innerSearch(conditions,child,&results)
	}
	return results
}

func innerSearch(conditions []string, pointer *Element, array *[]*Element){
	
	
	isMatch := true
	for _, cond := range conditions {
		condpair := strings.Split(cond,"=")
		key := strings.Replace(condpair[0]," ","",-1)
		values := strings.Replace(condpair[1], "\"", "", -1)
		values = strings.Replace(values, "'", "", -1)
		valuearray := strings.Fields(values)
		if key == "tag"{
			tagName := strings.Replace(valuearray[0]," ","",-1)
			if pointer.Tag != tagName{
				isMatch = false
			}
		}else{
			if val, ok := pointer.Attr[key]; ok {
				for _,valchild := range valuearray{
					htmlvals := strings.Split(val," ")
					foundval := false
					for _,htmlval := range htmlvals{
						if htmlval == valchild{
							foundval = true
						}
					}
					if !foundval{
						isMatch = false
					}
				}
			}else{
				isMatch = false
			}
		}
	}
	if isMatch{
		*array = append(*array, pointer)
	}
	
	for _, child := range pointer.Child {
		innerSearch(conditions,child,array)
	}
}
