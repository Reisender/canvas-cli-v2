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
# Standard mode - select one user at a time
canvas-cli users list [course-id]

# Multi-select mode - select multiple users
canvas-cli users list [course-id] --multi
```

In multi-select mode:
- Use up/down arrow keys to navigate
- Press space to select/deselect a user
- Press 'a' to select all users
- Press enter to show actions for the selected users

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

There are multiple ways to remove users from a course:

```bash
# Remove a single user by user ID
canvas-cli users remove [course-id] [user-id]

# Interactive removal - list users, select one, and choose "Remove"
canvas-cli users list [course-id]

# Bulk removal - list users in multi-select mode, select multiple users, and remove them
canvas-cli users list [course-id] --multi
```

The multi-select mode allows you to select and remove multiple users at once, which is much more efficient than removing them one by one.

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