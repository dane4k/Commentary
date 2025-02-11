## Запуск: 

## 1) .env:
```
DB_USER=postgres  # Пользователь PostgreSQL
DB_PASSWORD=admin  # Пароль пользователя
DB_NAME=commentary  # Название БД
```
## 2) config.yaml:
```
server:
  port: 8080

database:
  host: "db"
  port: 5432  # порт PostgreSQL
  user: "postgres"  # Пользователь PostgreSQL
  password: "admin"  # Пароль пользователя
  name: "commentary"  # Название БД
  store_in_db: true  # true - сохранение в БД, false - in-memory

logger:
  filename: "Commentary.log"
```
## 3) Из корня:
```
docker compose up --build
```

#### Сервисы: Инициализация БД, health-check -> применение миграций -> запуск приложения  


#
- Go
- PostgreSQL (squirrel)
- GraphQL (gqlgen)


#### Добавление отдельного query для получения комментариев по посту мне показалось оверкиллом, так как их можно получить через posts(limit: Int, offset: Int): [Post!]!

#### Добавление "точного" индекса id в in-memory мапы мне тоже показалось оверкиллом, так как фича удаления постов не требовалась, соответственно и точный индекс был бы бесполезен

#### Также посчитал, что sync.Map и разделение репозиториев по хранилищам в in-memory - оверкилл

```
scalar Time

type Post {
    id: ID!
    author: User!
    title: String!
    content: String!
    created: Time!
    commentable: Boolean!
    comments(limit: Int, offset: Int): [Comment!]!
}

type User {
    id: ID!
    username: String!
}

type Comment {
    id: ID!
    post: Post!
    author: User!
    content: String!
    created: Time!
    parent: Comment
    replies: [Comment!]!
}

input CreatePostInput {
    authorID: ID!
    title: String!
    content: String!
    commentable: Boolean!
}

input CreateCommentInput {
    authorID: ID!
    postID: ID!
    content: String!
    parent: ID
}

type Query {
    posts(limit: Int, offset: Int): [Post!]!

    post(postID: ID!, limit: Int, offset: Int): Post!
}

type Mutation {
    createUser(username: String!): User!

    createPost(input: CreatePostInput!): Post!

    toggleComments(postID: ID!): Post!

    createComment(input: CreateCommentInput!): Comment!
}

type Subscription {
    newComment(postID: ID!): Comment!
}
```
