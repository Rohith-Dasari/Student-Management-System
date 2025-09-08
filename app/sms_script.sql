-- SQLite
-- .tables
-- create table user(
-- UserID Text PRIMARY Key,
-- Name Text Not NUll,
-- Email Text Not Null,
-- Password Text Not Null,
-- Role Text Not Null Check(Role In ('faculty','student','admin'))DEFAULT 'faculty'
-- );


-- create table class(
-- ClassID Text PRIMARY Key,
-- Capacity Integer,
-- OccupiedBy Text
-- );


-- create table students(
-- StudentID text PRIMARY KEY,
-- Name text not null,
-- RollNumber Text UNIQUE NOT NULL,
-- ClassID Text not null,
-- semester integer not null,
-- FOREIGN key (ClassID) REFERENCES class(ClassID)
-- )


-- create table grades(
-- SubjectID text not null,
-- StudentID Text not null,
-- Grade INTEGER not NULL,
-- semester integer not null,
-- PRIMARY KEY(SubjectID,StudentID),
-- FOREIGN Key(SubjectID) REFERENCES subject(SubjectID),
-- FOREIGN Key(StudentID) REFERENCES students(StudentID)
-- );


-- insert into user VALUES(
-- 'dafceae6-fd26-4e40-8609-ed4bdb80ddf6',
-- 'Rohith',
-- 'admin@gmail.com',
-- '$2a$10$G01O6/zbOQ3cyZeAJ9WHSuPbvZ90AlKesKhDfq56359e7eyRGnbyS',
-- 'admin'
-- )
-- insert into class (ClassID,Capacity) values(
-- 'f0846391-ebc2-48dd-b1b4-33c791096ab6',
-- 60
-- );


-- insert into subject VALUES(
-- '1250d94e-38a5-4734-be0d-048c6e6d51c1',
-- 'Engineering Mathematics'
-- )
-- insert into subject VALUES(
-- '4f4c7579-7aa3-47cd-8e2f-c2580aa3d8f6',
-- 'Physics'
-- )
