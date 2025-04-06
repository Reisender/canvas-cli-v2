# Canvas CLI

A command line interface for interacting with the Canvas LMS API, built with Go and [Charm](https://charm.sh) libraries.

## Features

- List and view courses
- List and view assignments
- Manage course users and enrollments
- Easy configuration management
- Beautiful terminal UI using Charm libraries

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/Reisender/canvas-cli-v2.git
cd canvas-cli-v2

# Build the binary
go build -o canvas-cli ./cmd/canvas-cli

# Move to a directory in your PATH (optional)
sudo mv canvas-cli /usr/local/bin/
```

## Usage

### First-time Setup

Before using Canvas CLI, you need to configure it with your Canvas API key:

```bash
canvas-cli config
```

You'll be prompted to enter:
- Canvas API base URL (defaults to https://canvas.instructure.com/api/v1)
- Your Canvas API key

You can also set these values individually:

```bash
canvas-cli config set api_key your-api-key
canvas-cli config set base_url https://your-institution.instructure.com/api/v1
```

### View Your Configuration

```bash
canvas-cli config get
```

### List Your Courses

```bash
canvas-cli courses list
```

### View Course Assignments

```bash
canvas-cli assignments list [course-id]
```

### Managing Users in a Course

#### List Users in a Course

```bash
canvas-cli users list [course-id]
```

#### View User Details

```bash
canvas-cli users view [user-id]
```

#### List Enrollments in a Course

```bash
canvas-cli users enrollments list [course-id]
```

#### Add a User to a Course

```bash
# Enroll a user as a student
canvas-cli users enrollments add [course-id] [user-id]

# Enroll a user as a teacher
canvas-cli users enrollments add [course-id] [user-id] --type TeacherEnrollment

# Enroll a user and send notification
canvas-cli users enrollments add [course-id] [user-id] --notify
```

Available enrollment types:
- StudentEnrollment (default)
- TeacherEnrollment
- TaEnrollment
- DesignerEnrollment
- ObserverEnrollment

#### Remove a User from a Course

There are two ways to remove a user from a course:

```bash
# Remove by user ID (easier)
canvas-cli users remove [course-id] [user-id]

# Remove by enrollment ID
canvas-cli users enrollments remove [course-id] [enrollment-id]
```

Using the `users remove` command is recommended as it only requires the user ID, which is easier to find. The enrollment ID can be found by using the `users enrollments list` command if needed.

## Development

### Requirements

- Go 1.24 or higher
- Dependencies:
  - github.com/charmbracelet/bubbles
  - github.com/charmbracelet/bubbletea
  - github.com/charmbracelet/lipgloss
  - github.com/spf13/cobra
  - github.com/spf13/viper

### Building from Source

```bash
go build -o canvas-cli ./cmd/canvas-cli
```

## License

MIT 