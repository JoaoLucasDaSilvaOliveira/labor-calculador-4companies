# LaborCalculador4Companies
![Go|111](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Fyne|102](https://img.shields.io/badge/Fyne-00D1B2?style=for-the-badge&logo=fyne&logoColor=white)
![Go Version|211](https://img.shields.io/github/go-mod/go-version/JoaoLucasDaSilvaOliveira/learning-fyne?style=for-the-badge)
This file contains the processes and decisions for creation of screens and designs of the LaborCalculador4Companies system, a desktop application for companies to generate labor calculations with professional expertise.
## Toolkit Used

### Front End
#### Prototyping
- Excalidraw for simple assemble ideas and create a basic/rudimentary view. Also the it's gonna be used for brainstorming and review the system architecture/design.
- Figma, on a more advanced state, for create interactions, real application flow and some stuff like colors, boxes and widget behaviors and other things like this.
#### Coding
Fyne Framework. Go lang GUI provider.
### Back End
#### Language
Go. Patterns: DDD, Hexagonal, Clean Architecture and Factory.
#### Testing
Native language test pattern.
#### For Infrastructure
SQLite database for local access and MVP. Self-hosted PostgreSQL database when web app is released.
### Overall
#### Build
- Into linux environment: `go build` with `GOOS=linux && GOARCH=amd64`.
- Into windows environment: `fyne-crosss` builder, with the same amd64 architecture.
- Docker strict when web app is released.
#### Version control and download
Git and GitHub. Also gonna use GitHub Releases for downloading and updating the desktop application.
