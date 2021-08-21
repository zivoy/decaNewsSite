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

const tagPath = adminBasePath + "/tag_list"

func getAvailableTags() map[string]Tag {
	items := tagListCache.get("tags", func(string) interface{} {
		ref := dataBase.NewRef(tagPath)
		var data []Tag
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

// get list of tags from comma seperated list
func getTagsFromString(s string) []Tag {
	tags := getAvailableTags()
	list := strings.Split(s, ",")
	tagList := make([]Tag, len(list))
	tagsPut := map[string]bool{}

	for i, v := range list {
		id := strings.Trim(strings.ToLower(v), " ")
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
	tagMap := getAvailableTags()
	var tags = make([]Tag, 0)
	for _, v := range tagMap {
		tags = append(tags, v)
	}
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

func compareTagList(a, b []Tag) bool {
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

//todo cleanup unused tags
