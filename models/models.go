package models

import (
    "database/sql"
    "time"
)


type Building struct {
    ID         int    `json:"building_id"`
    Name       string `json:"building_name"`
    GenderType string `json:"gender_type"`
}

func GetAllBuildings(db *sql.DB) ([]Building, error) {
    rows, err := db.Query("SELECT building_id, building_name, gender_type FROM building")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var buildings []Building
    for rows.Next() {
        var b Building
        if err := rows.Scan(&b.ID, &b.Name, &b.GenderType); err != nil {
            return nil, err
        }
        buildings = append(buildings, b)
    }
    return buildings, nil
}


func CreateBuilding(db *sql.DB, name, genderType string) error {
    _, err := db.Exec("INSERT INTO building (building_name, gender_type) VALUES (?, ?)", name, genderType)
    return err
}


type Floor struct {
    ID         int `json:"floor_id"`
    Number     int `json:"floor_number"`
    BuildingID int `json:"building_id"`
}


func GetAllFloors(db *sql.DB) ([]Floor, error) {
    rows, err := db.Query("SELECT floor_id, floor_number, building_id FROM floor")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var floors []Floor
    for rows.Next() {
        var f Floor
        if err := rows.Scan(&f.ID, &f.Number, &f.BuildingID); err != nil {
            return nil, err
        }
        floors = append(floors, f)
    }
    return floors, nil
}

func CreateFloor(db *sql.DB, number, buildingID int) error {
    _, err := db.Exec("INSERT INTO floor (floor_number, building_id) VALUES (?, ?)", number, buildingID)
    return err
}

type Room struct {
    ID        int    `json:"room_id"`
    Number    int    `json:"room_number"`
    Capacity  int    `json:"capacity"`
    Amenities string `json:"amenities"`
    FloorID   int    `json:"floor_id"`
}


func GetAllRooms(db *sql.DB) ([]Room, error) {
    rows, err := db.Query("SELECT room_id, room_number, capacity, amenities, floor_id FROM room")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var rooms []Room
    for rows.Next() {
        var r Room
        var amenities sql.NullString
        if err := rows.Scan(&r.ID, &r.Number, &r.Capacity, &amenities, &r.FloorID); err != nil {
            return nil, err
        }
        r.Amenities = amenities.String
        rooms = append(rooms, r)
    }
    return rooms, nil
}

func CreateRoom(db *sql.DB, number, capacity, floorID int, amenities string) error {
    _, err := db.Exec("INSERT INTO room (room_number, capacity, amenities, floor_id) VALUES (?, ?, ?, ?)", number, capacity, amenities, floorID)
    return err
}


func UpdateRoomAmenities(db *sql.DB, roomID int, amenities string) error {
    _, err := db.Exec("CALL update_room_amenities(?, ?)", roomID, amenities)
    return err
}


type Department struct {
    ID   int    `json:"department_id"`
    Name string `json:"department_name"`
}

func GetAllDepartments(db *sql.DB) ([]Department, error) {
    rows, err := db.Query("SELECT department_id, department_name FROM department")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var departments []Department
    for rows.Next() {
        var d Department
        if err := rows.Scan(&d.ID, &d.Name); err != nil {
            return nil, err
        }
        departments = append(departments, d)
    }
    return departments, nil
}

func CreateDepartment(db *sql.DB, name string) error {
    _, err := db.Exec("INSERT INTO department (department_name) VALUES (?)", name)
    return err
}


type Student struct {
    ID           int    `json:"student_id"`
    FirstName    string `json:"first_name"`
    LastName     string `json:"last_name"`
    Gender       string `json:"gender"`
    RoomID       int    `json:"room_id"`
    DepartmentID int    `json:"department_id"`
    BuildingID   int    `json:"building_id"`
}

func GetAllStudents(db *sql.DB) ([]Student, error) {
    rows, err := db.Query("SELECT student_id, first_name, last_name, gender, room_id, department_id, building_id FROM student")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var students []Student
    for rows.Next() {
        var s Student
        if err := rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Gender, &s.RoomID, &s.DepartmentID, &s.BuildingID); err != nil {
            return nil, err
        }
        students = append(students, s)
    }
    return students, nil
}


func CreateStudent(db *sql.DB, firstName, lastName, gender string, roomID, departmentID, buildingID int) error {
    _, err := db.Exec("CALL assign_student_to_room(?, ?, ?, ?, ?, ?)", firstName, lastName, gender, roomID, departmentID, buildingID)
    return err
}


func UpdateStudent(db *sql.DB, id int, firstName, lastName, gender string, roomID, departmentID, buildingID int) error {
    _, err := db.Exec("UPDATE student SET first_name = ?, last_name = ?, gender = ?, room_id = ?, department_id = ?, building_id = ? WHERE student_id = ?",
        firstName, lastName, gender, roomID, departmentID, buildingID, id)
    return err
}

func DeleteStudent(db *sql.DB, id int) error {
    _, err := db.Exec("DELETE FROM student WHERE student_id = ?", id)
    return err
}


type AuditLog struct {
    ID         int       `json:"log_id"`
    TableName  string    `json:"table_name"`
    Operation  string    `json:"operation"`
    RecordID   int       `json:"record_id"`
    ChangeTime time.Time `json:"change_time"`
    ChangedBy  string    `json:"changed_by"`
    Details    string    `json:"details"`
}


func GetAllAuditLogs(db *sql.DB) ([]AuditLog, error) {
    rows, err := db.Query("SELECT log_id, table_name, operation, record_id, change_time, changed_by, details FROM audit_log")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var logs []AuditLog
    for rows.Next() {
        var l AuditLog
        if err := rows.Scan(&l.ID, &l.TableName, &l.Operation, &l.RecordID, &l.ChangeTime, &l.ChangedBy, &l.Details); err != nil {
            return nil, err
        }
        logs = append(logs, l)
    }
    return logs, nil
}