# Todo App
To understand and work with go, htmx and sqlite3 to make a simple todo app. The app is not going to be online ever. The idea is to have a CLI interface and Web interface for GUI.

## Goals
- [x] Create a Server which can server html pages.
- [x] Add a styling system.
- [x] Add htmx for interaction handling.
- [ ] Connect to a local sqlite3 db.
- [ ] Make a CLI interface.
- [ ] Organise the Project in multiple module.


## Learnings
- Go template is good enough to solve this problem. The type safety of templ is good, But it is not bad right now.
- TailwindCss is not bad. I enjoyed using it as utility css. Though I did add 
one custom css in. You can see that in `input.css` file. It handle's my logic of
having different style when the task is marked done.

## Run commands
```bash
go run main.go
```

## Development commands
```bash
npm run tailwind
go run main.go
```

## Current State
![Todo Web View](https://github.com/ronakmehtav/httpGo/assets/31774137/6814b0f1-e811-4fe5-95cf-b6b882798d09)

