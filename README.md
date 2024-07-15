# Technical Challenge from fliQt
interviewee: [BoHsiang Huang - calvin.peak@finalbuild.zip](mailto:calvin.peak@finalbuild.zip)

# Designing a HR system API
Since this project is built for an interview, I wanna ensure that we focus on the most core part, which is demonstrating my technical skills. I have marked some of the functional requirements below and tried to make the project smaller. Perhaps we can discuss the items that were deleted at some point.

## Functional requirements
- There are 3 different roles of characters using this system: candidate, interviewer, HR.
- Candidates can search for jobs and filter by title, job description, salary range and job type.
- Candidates can check their application status.
- ~~HR can use the system to manage schedules and status.~~
- HR can manage applications from candidates.
- A job can be created/updated, or closed by HR.
- ~~The system can organize interview schedule.~~

## Non-functional requirements
- Since resumes contain highly confidential data, interviewers and HR must pass the 2FA before downloading them.
- ~~Analysing resumes and making a score for each candidates.~~
- Records every operation for tracking purposes.

I also simplified user authentication for users. You can easily use the `X-FLIQT-USER` header to interact with API as a specific user.
```sh
$ curl -H 'X-FLIQT-USER: [candidate 1's id]' http://localhost:8080/api/applications
```

# How to Run the Code
## Environment Variables
| Name | Description | Default |
|------|-------------|---------|
|`DB_HOST`| DB host (This project only supports MySQL) | `localhost`|
|`DB_PORT`| DB port (This project only supports MySQL) | `3306`|
|`DB_USER`| DB username | `root`|
|`DB_PASSWORD`| DB password | `password` |
|`DB_TIMEZONE`| DB timezone | `Etc/UTC` |
|`DB_NAME`| DB name | `fliqt_test`|
|`DB_MAX_IDLE`| DB maximum connections in idle connection pool | 2 |
|`DB_MAX_CONN`| DB maximum number of open connections | 10 |
|`DB_MAX_LIFE`| DB maximum amount of time  amount of time a connection may be reused | 2 |
|`DEBUG`| Config for debugging mode | `false` |
|`PRETTY_LOG`| Easier to read the log | `false` |
|`TRACER_ENDPOINT`| The tracer collector's endpoint | `localhost:4317` |


