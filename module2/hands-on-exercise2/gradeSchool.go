package main

import (
    "fmt"
    "sort"
)

type Grade struct {
    Grade int
    Students []string
}

type School map[int]*Grade

func New() *School {
    return &School{}
}

func (s *School) Add(name string,grade int){
    if _,ok := (*s)[grade] ; !ok {
        (*s)[grade] = &Grade{
            Grade : grade,
            Students : []string{name},
            }
        } else if !s.isStudentPresent(name,grade) {
            (*s)[grade].Students = append((*s)[grade].Students,name)
            sort.Strings((*s)[grade].Students)
        }
}

func (s *School) Enrollement() (output []Grade) {
    for _,grade := range *s {
        output = append(output,* grade)
    }
    sort.Slice(output,func(i,j int) bool {
        return output[i].Grade < output[j].Grade
    })
    return
}

func (s School) isStudentPresent (name string,grade int) bool {
	for _,ele := range s[grade].Students {
		if ele == name {
			return true
		}
	}
	return false
}


func gradeSchool() {
    school := New()
    
    school.Add("Vashista",10)
    school.Add("Pranaw",10)
    school.Add("Resha",10)
    school.Add("Prajakta",6)
    school.Add("Ashish",6)
	school.Add("Ashish",6)
    
    fmt.Println(school.Enrollement())
}