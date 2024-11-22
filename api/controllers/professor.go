package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/api/common/log"
	"github.com/UTDNebula/nebula-api/api/configs"
	"github.com/UTDNebula/nebula-api/api/responses"
	"github.com/UTDNebula/nebula-api/api/schema"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var professorCollection *mongo.Collection = configs.GetCollection("professors")

// @Id professorSearch
// @Router /professor [get]
// @Description "Returns all professors matching the query's string-typed key-value pairs"
// @Produce json
// @Param first_name query string false "The professor's first name"
// @Param last_name query string false "The professor's last name"
// @Param titles query string false "One of the professor's title"
// @Param email query string false "The professor's email address"
// @Param phone_number query string false "The professor's phone number"
// @Param office.building query string false "The building of the location of the professor's office"
// @Param office.room query string false "The room of the location of the professor's office"
// @Param office.map_uri query string false "A hyperlink to the UTD room locator of the professor's office"
// @Param profile_uri query string false "A hyperlink pointing to the professor's official university profile"
// @Param image_uri query string false "A link to the image used for the professor on the professor's official university profile"
// @Param office_hours.start_date query string false "The start date of one of the office hours meetings of the professor"
// @Param office_hours.end_date query string false "The end date of one of the office hours meetings of the professor"
// @Param office_hours.meeting_days query string false "One of the days that one of the office hours meetings of the professor"
// @Param office_hours.start_time query string false "The time one of the office hours meetings of the professor starts"
// @Param office_hours.end_time query string false "The time one of the office hours meetings of the professor ends"
// @Param office_hours.modality query string false "The modality of one of the office hours meetings of the professor"
// @Param office_hours.location.building query string false "The building of one of the office hours meetings of the professor"
// @Param office_hours.location.room query string false "The room of one of the office hours meetings of the professor"
// @Param office_hours.location.map_uri query string false "A hyperlink to the UTD room locator of one of the office hours meetings of the professor"
// @Param sections query string false "The _id of one of the sections the professor teaches"
// @Success 200 {array} schema.Professor "A list of professors"
func ProfessorSearch(c *gin.Context) {
	//name := c.Query("name")            // value of specific query parameter: string
	//queryParams := c.Request.URL.Query() // map of all query params: map[string][]string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var professors []schema.Professor

	defer cancel()

	// build query key value pairs (only one value per key)
	query, err := schema.FilterQuery[schema.Professor](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "schema validation error", Data: err.Error()})
		return
	}

	optionLimit, err := configs.GetOptionLimit(&query, c)
	if err != nil {
		log.WriteErrorWithMsg(err, log.OffsetNotTypeInteger)
		c.JSON(http.StatusConflict, responses.ErrorResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
		return
	}

	// get cursor for query results
	cursor, err := professorCollection.Find(ctx, query, optionLimit)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// retrieve and parse all valid documents
	if err = cursor.All(ctx, &professors); err != nil {
		log.WritePanic(err)
		panic(err)
	}

	// return result
	c.JSON(http.StatusOK, responses.MultiProfessorResponse{Status: http.StatusOK, Message: "success", Data: professors})
}

// @Id professorById
// @Router /professor/{id} [get]
// @Description "Returns the professor with given ID"
// @Produce json
// @Param id path string true "ID of the professor to get"
// @Success 200 {object} schema.Professor "A professor"
func ProfessorById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	professorId := c.Param("id")

	var professor schema.Professor

	defer cancel()

	// parse object id from id parameter
	objId, err := primitive.ObjectIDFromHex(professorId)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
		return
	}

	// find and parse matching professor
	err = professorCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&professor)
	if err != nil {
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// return result
	c.JSON(http.StatusOK, responses.SingleProfessorResponse{Status: http.StatusOK, Message: "success", Data: professor})
}

func ProfessorAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	var professors []schema.Professor

	defer cancel()

	cursor, err := professorCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}

	// retrieve and parse all valid documents
	if err = cursor.All(ctx, &professors); err != nil {
		panic(err)
	}

	// return result
	c.JSON(http.StatusOK, responses.MultiProfessorResponse{Status: http.StatusOK, Message: "success", Data: professors})
}

// @Id professorCourseSearch
// @Router /professor/courses [get]
// @Description "Returns all of the courses of all the professors matching the query's string-typed key-value pairs"
// @Produce json
// @Param first_name query string false "The professor's first name"
// @Param last_name query string false "The professor's last name"
// @Param titles query string false "One of the professor's title"
// @Param email query string false "The professor's email address"
// @Param phone_number query string false "The professor's phone number"
// @Param office.building query string false "The building of the location of the professor's office"
// @Param office.room query string false "The room of the location of the professor's office"
// @Param office.map_uri query string false "A hyperlink to the UTD room locator of the professor's office"
// @Param profile_uri query string false "A hyperlink pointing to the professor's official university profile"
// @Param image_uri query string false "A link to the image used for the professor on the professor's official university profile"
// @Param office_hours.start_date query string false "The start date of one of the office hours meetings of the professor"
// @Param office_hours.end_date query string false "The end date of one of the office hours meetings of the professor"
// @Param office_hours.meeting_days query string false "One of the days that one of the office hours meetings of the professor"
// @Param office_hours.start_time query string false "The time one of the office hours meetings of the professor starts"
// @Param office_hours.end_time query string false "The time one of the office hours meetings of the professor ends"
// @Param office_hours.modality query string false "The modality of one of the office hours meetings of the professor"
// @Param office_hours.location.building query string false "The building of one of the office hours meetings of the professor"
// @Param office_hours.location.room query string false "The room of one of the office hours meetings of the professor"
// @Param office_hours.location.map_uri query string false "A hyperlink to the UTD room locator of one of the office hours meetings of the professor"
// @Param sections query string false "The _id of one of the sections the professor teaches"
// @Success 200 {array} schema.Course "A list of Courses"
func ProfessorCourseSearch() gin.HandlerFunc {
	// Wrapper of professorCourse() with flag of Search
	return func(c *gin.Context) {
		professorCourse("Search", c)
	}
}

// @Id professorCourseById
// @Router /professor/{id}/courses [get]
// @Description "Returns all the courses taught by the professor with given ID"
// @Produce json
// @Param id path string true "ID of the professor to get"
// @Success 200 {array} schema.Course "A list of courses"
func ProfessorCourseById() gin.HandlerFunc {
	// Essentially wrapper of professorCourse() with flag of ById
	return func(c *gin.Context) {
		professorCourse("ById", c)
	}
}

// Get all of the courses of the professors depending on the type of flag
func professorCourse(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var professorCourses []schema.Course // array of courses of the professors (or single professor with Id)
	var professorQuery bson.M            // query filter the professor
	var err error

	defer cancel()

	// determine the professor's query
	if professorQuery, err = getProfessorQuery(flag, c); err != nil {
		return // if there's an error, the response will have already been thrown to the consumer, halt the funcion here
	}

	// determine the offset and limit for pagination stage
	// and delete "offset" field in professorQuery
	offset, limit, err := configs.GetAggregateLimit(&professorQuery, c)
	if err != nil {
		log.WriteErrorWithMsg(err, log.OffsetNotTypeInteger)
		c.JSON(http.StatusConflict, responses.ErrorResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
		return
	}

	// Pipeline to query the courses from the filtered professors (or a single professor)
	professorCoursePipeline := mongo.Pipeline{
		// filter the professors
		bson.D{{Key: "$match", Value: professorQuery}},

		// paginate the professors before pulling the courses from those professor
		bson.D{{Key: "$skip", Value: offset}}, // skip to the specified offset
		bson.D{{Key: "$limit", Value: limit}}, // limit to the specified number of professors

		// lookup the array of sections from sections collection
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sections"},
			{Key: "localField", Value: "sections"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sections"},
		}}},

		// project the courses referenced by each section in the array
		bson.D{{Key: "$project", Value: bson.D{{Key: "courses", Value: "$sections.course_reference"}}}},

		// lookup the array of courses from coures collection
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "courses"},
			{Key: "localField", Value: "courses"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "courses"},
		}}},

		// unwind the courses
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$courses"},
			{Key: "preserveNullAndEmptyArrays", Value: false}, // to avoid the professor documents that can't be replaced
		}}},

		// replace the combination of ids and courses with the courses entirely
		bson.D{{Key: "$replaceWith", Value: "$courses"}},
	}

	// Perform aggreration on the pipeline
	cursor, err := professorCollection.Aggregate(ctx, professorCoursePipeline)
	if err != nil {
		// return the error with there's something wrong with the aggregation
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}
	// Parse the array of courses from these professors
	if err = cursor.All(ctx, &professorCourses); err != nil {
		log.WritePanic(err)
		panic(err)
	}
	c.JSON(http.StatusOK, responses.MultiCourseResponse{Status: http.StatusOK, Message: "success", Data: professorCourses})
}

// @Id professorSectionSearch
// @Router /professor/sections [get]
// @Description "Returns all of the sections of all the professors matching the query's string-typed key-value pairs"
// @Produce json
// @Param first_name query string false "The professor's first name"
// @Param last_name query string false "The professor's last name"
// @Param titles query string false "One of the professor's title"
// @Param email query string false "The professor's email address"
// @Param phone_number query string false "The professor's phone number"
// @Param office.building query string false "The building of the location of the professor's office"
// @Param office.room query string false "The room of the location of the professor's office"
// @Param office.map_uri query string false "A hyperlink to the UTD room locator of the professor's office"
// @Param profile_uri query string false "A hyperlink pointing to the professor's official university profile"
// @Param image_uri query string false "A link to the image used for the professor on the professor's official university profile"
// @Param office_hours.start_date query string false "The start date of one of the office hours meetings of the professor"
// @Param office_hours.end_date query string false "The end date of one of the office hours meetings of the professor"
// @Param office_hours.meeting_days query string false "One of the days that one of the office hours meetings of the professor"
// @Param office_hours.start_time query string false "The time one of the office hours meetings of the professor starts"
// @Param office_hours.end_time query string false "The time one of the office hours meetings of the professor ends"
// @Param office_hours.modality query string false "The modality of one of the office hours meetings of the professor"
// @Param office_hours.location.building query string false "The building of one of the office hours meetings of the professor"
// @Param office_hours.location.room query string false "The room of one of the office hours meetings of the professor"
// @Param office_hours.location.map_uri query string false "A hyperlink to the UTD room locator of one of the office hours meetings of the professor"
// @Param sections query string false "The _id of one of the sections the professor teaches"
// @Success 200 {array} schema.Section "A list of Sections"
func ProfessorSectionSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		professorSection("Search", c)
	}
}

// @Id professorSectionById
// @Router /professor/{id}/sections [get]
// @Description "Returns all the sections taught by the professor with given ID"
// @Produce json
// @Param id path string true "ID of the professor to get"
// @Success 200 {array} schema.Section "A list of sections"
func ProfessorSectionById() gin.HandlerFunc {
	return func(c *gin.Context) {
		professorSection("ById", c)
	}
}

// Get all of the sections of the professors depending on the type of flag
func professorSection(flag string, c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var professorSections []schema.Section // array of sections of the professors (or single professor with Id)
	var professorQuery bson.M              // query filter the professor
	var err error

	defer cancel()

	// determine the professor's query
	if professorQuery, err = getProfessorQuery(flag, c); err != nil {
		return // if there's an error, the response will have already been thrown to the consumer, halt the funcion here
	}

	// determine the offset and limit for pagination stage
	// and delete "offset" field in professorQuery
	offset, limit, err := configs.GetAggregateLimit(&professorQuery, c)
	if err != nil {
		log.WriteErrorWithMsg(err, log.OffsetNotTypeInteger)
		c.JSON(http.StatusConflict, responses.ErrorResponse{Status: http.StatusConflict, Message: "Error offset is not type integer", Data: err.Error()})
		return
	}

	// Pipeline to query the courses from the filtered professors (or a single professor)
	professorSectionPipeline := mongo.Pipeline{
		// filter the professors
		bson.D{{Key: "$match", Value: professorQuery}},

		// paginate the professors before pulling the courses from those professor
		bson.D{{Key: "$skip", Value: offset}}, // skip to the specified offset
		bson.D{{Key: "$limit", Value: limit}}, // limit to the specified number of professors

		// lookup the array of sections from sections collection
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "sections"},
			{Key: "localField", Value: "sections"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sections"},
		}}},

		// project the sections
		bson.D{{Key: "$project", Value: bson.D{{Key: "sections", Value: "$sections"}}}},

		// unwind the sections
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$sections"},
			{Key: "preserveNullAndEmptyArrays", Value: false}, // to avoid the professor documents that can't be replaced
		}}},

		// replace the combination of ids and sections with the sections entirely
		bson.D{{Key: "$replaceWith", Value: "$sections"}},
	}

	// Perform aggreration on the pipeline
	cursor, err := professorCollection.Aggregate(ctx, professorSectionPipeline)
	if err != nil {
		// return the error with there's something wrong with the aggregation
		log.WriteError(err)
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		return
	}
	// Parse the array of sections from these professors
	if err = cursor.All(ctx, &professorSections); err != nil {
		log.WritePanic(err)
		panic(err)
	}
	c.JSON(http.StatusOK, responses.MultiSectionResponse{Status: http.StatusOK, Message: "success", Data: professorSections})
}

// function to determine the query of the professor based on the parameters passed from context
// avoid redundancy in the code
// if there's an error, throw an error response back to the API consumer and return only the error
func getProfessorQuery(flag string, c *gin.Context) (bson.M, error) {
	var professorQuery bson.M
	var err error

	if flag == "Search" { // if the flag is Search, filter professors based on query parameters
		// build the key-value pairs of query parameters
		professorQuery, err = schema.FilterQuery[schema.Professor](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "schema validation error", Data: err.Error()})
			return nil, err // return only the error
		}
	} else if flag == "ById" { // if the flag is ById, filter that single professor based on their _id
		// parse the ObjectId
		professorId := c.Param("id")
		professorObjId, convertIdErr := primitive.ObjectIDFromHex(professorId)
		if convertIdErr != nil {
			log.WriteError(convertIdErr)
			c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest, Message: "id conversion error", Data: convertIdErr.Error()})
			return nil, convertIdErr
		}
		professorQuery = bson.M{"_id": professorObjId}
	} else {
		err = errors.New("broken endpoint")
		// something wrong that messed up the server
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, Message: "endpoint error", Data: err.Error()})
		return nil, err
	}
	return professorQuery, err
}
