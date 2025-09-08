# Live Portfolio

A dynamic portfolio site powered by a **C/GO backend** (In deciding stage) and **React/CSS frontend**.

## Features

- React UI that fetches live portfolio data
- C/GO backend providing JSON API (mock data, can be extended)
- Easy to run locally

## Stack

- Backend: C/Go (REST API, mock portfolio data)
- Frontend: React + CSS

## Running Locally


### 1. Start React Frontend
cd frontend-react
npm install
npm start

By default, React runs on `localhost:3000` and C backend on `localhost:8080`.

## Notes

- Idea is to eventually scrape my LinkedIn in order to autofill the information to be put on the portfolio  -> should be feasible?
-   - May also be worth scraping resumes?
- Backend is still in a preliminary stage, and I'm currently unsure which language to choose. Options include C & Go. ---> I do not think running C backend alone is feasible but may just shift fully to Go. (This may be the reason why C is not used in web development)
- Frontend server works albeit it just says 'Loading ...' currently.
- Lets try svelte/more js frontend? -- DENIED --
- also need to figure out if hosting is possible including backend instead of boring frontend that isnt adaptable since the backend isnt attached ----> Should be feasible with Go backend
