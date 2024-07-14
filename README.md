# Technical Challenge from fliQt
interviewee: [BoHsiang Huang - calvin.peak@finalbuild.zip](mailto:calvin.peak@finalbuild.zip)

# Designing a HR system API

## Functional requirements
- There are 3 different roles of characters using this system: candidate, interviewer, HR.
- Candidates can search for jobs and filter by title, job description, salary range and job type.
- HR can use the system to manage schedules and status.
- A job can be created/updated, or closed by HR.
- The system can organize interview schedule.

## Non-functional requirements
- Since resumes contain highly confidential data, interviewers and HR must pass the 2FA before downloading them.
- Analysing resumes and making a score for each candidates.
- Records every operation for tracking purposes.
