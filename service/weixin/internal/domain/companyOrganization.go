package domain

type OrganizationCourses []OrganizationCourse

type OrganizationCourseModule struct {
	CourseModuleName    string `json:"courseModuleName"`
	CourseModuleContent string `json:"courseModuleContent"`
}

type OrganizationCourse struct {
	CourseName          string                     `json:"courseName"`
	CourseSubName       string                     `json:"courseSubName"`
	CoursePrice         float64                    `json:"coursePrice"`
	CourseDuration      uint64                     `json:"courseDuration"`
	CourseOriginalPrice float64                    `json:"courseOriginalPrice"`
	CourseLevel         uint8                      `json:"courseLevel"`
	CourseModules       []OrganizationCourseModule `json:"courseModule"`
}

func (o OrganizationCourses) Len() int {
	return len(o)
}

func (o OrganizationCourses) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o OrganizationCourses) Less(i, j int) bool {
	return o[i].CoursePrice < o[j].CoursePrice
}
