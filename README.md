# Live Portfolio

A dynamic portfolio site powered by a **C backend** and **React/CSS frontend**.

## Features

- React UI that fetches live portfolio data
- C backend providing JSON API (mock data, can be extended)
- Easy to run locally

## Stack

- Backend: C (REST API, mock portfolio data)
- Frontend: React + CSS

## Running Locally

### 1. Start C Backend

cd backend-c
make
./server

### 2. Start React Frontend
cd frontend-react
npm install
npm start

By default, React runs on `localhost:3000` and C backend on `localhost:8080`.

## Notes

- **LinkedIn integration not implemented here**; see comments for where to add it!
- To extend the backend, edit `server.c` or use a database.
