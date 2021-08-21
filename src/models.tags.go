package main

import (
	"log"
	"sort"
	"strings"
)

type Tag struct {
	Name  string     `json:"name"`
	Color BulmaColor `json:"color_mask"`
	id    string
}

type tagList []Tag

const tagPath = adminBasePath + "/tag_list"

func getAvailableTags() map[string]Tag {
	items := tagListCache.get("tags", func(string) interface{} {
		ref := dataBase.NewRef(tagPath)
		var data tagList
		if err := ref.Get(ctx, &data); err != nil && debug {
			log.Println(err)
			return nil
		}

		tags := map[string]Tag{}
		for _, v := range data {
			v.id = strings.ToLower(v.Name)
			tags[v.id] = v
		}

		return tags
	})
	return items.(map[string]Tag)
}

func getTagList() tagList {
	tagMap := getAvailableTags()
	var tags = make(tagList, 0)
	for _, v := range tagMap {
		tags = append(tags, v)
	}
	return tags
}

// get list of tags from comma seperated list
func getTagsFromString(s string) tagList {
	if s == "" {
		return make(tagList, 0)
	}

	tags := getAvailableTags()
	list := strings.Split(s, ",")
	tagList := make(tagList, len(list))
	tagsPut := map[string]bool{}

	for i, v := range list {
		id := strings.Trim(strings.ToLower(v), " ")
		if id == "" {
			continue
		}
		if _, ok := tagsPut[id]; ok {
			continue
		}

		if tag, ok := tags[id]; ok {
			tagList[i] = tag
		} else {
			tagList[i] = Tag{
				Name:  strings.Trim(v, " "),
				Color: Info | Lighter,
				id:    id,
			}
			newTag(tagList[i])
		}
		tagsPut[id] = true
	}
	return tagList
}

func newTag(a Tag) {
	tags := getTagList()
	tags = append(tags, a)
	err := setEntry(dataBase, tagPath, tags)
	if err != nil {
		log.Println(err)
		return
	}

	clearData := cacheAction{
		CacheListId: tagListCache.id,
		ItemId:      "tags",
		ActionType:  clearList,
	}
	sendAction(clearData)
	tagListCache.clear()
}

func compareTagList(a, b tagList) bool {
	aIds := make([]string, len(a))
	bIds := make([]string, len(b))
	sort.Strings(aIds)
	sort.Strings(bIds)
	listA := strings.Join(aIds, ",")
	listB := strings.Join(bIds, ",")
	return listA == listB
}

func getTagFromID(id string) Tag {
	if tag, ok := getAvailableTags()[id]; ok {
		return tag
	}
	return Tag{Name: id, Color: Danger}
}

func (t tagList) String() string {
	names := make([]string, len(t))
	for i, v := range t {
		names[i] = v.Name
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}

//todo cleanup unused tags
