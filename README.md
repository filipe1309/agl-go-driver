<a name="readme-top"></a>

# <p align="center">agl-go-driver</p>

<p align="center">
    <img src="https://img.shields.io/badge/Code-Go-informational?style=flat-square&logo=go&color=00ADD8" alt="Go" />
</p>

## 💬 About

This project was developed following Aprenda Golang's "[Imersão Backend](https://aprendagolang.com.br/courses/imersao-backend)" course.

Notes taken during the course are in the [notes](notes.md) file.
Architecture:

## :art: Architecture

> [!NOTE]
> 1. The following diagram was generated using [Mermaid](https://mermaid-js.github.io/mermaid/).  
> 2. It is a C4 model diagram at the container level, see more at [C4 Model](https://c4model.com/).


```mermaid
  C4Container
    title Personal Go Drive
    Person(customerA, "Customer A", "A customer of the drive API")

    Enterprise_Boundary(b0, "System") {  
      
      Container(web_interface, "Web Application", "Interface that allows consumers to access Drive API via Web")
      
      Container(api, "API", "Allows customers upload/download files from Drive")

      Container(cli_interface, "CLI", "Interface that allows consumers to access Drive API via terminal")

      Container(worker, "Worker", "Compress files and remove raw files")
      
      ContainerDb(db, "Storage", "Stores the links of the uploaded files")
      
      ContainerQueue(queue, "Queue")
    }

    System(bucket_raw, "Bucket Raw", "Save raw files")
    System(bucket_compress, "Bucket Compress", "Save compress files")

    BiRel(customerA, web_interface, "HTTP")
    BiRel(customerA, cli_interface, "Terminal")
    Rel_D(cli_interface, api, "Save/Get link")
    Rel_D(web_interface, api, "Save/Get link")
    Rel_D(api, db, "Save/Get link")
    Rel_D(api, queue, "Save/Get link")
    Rel_D(api, bucket_raw, "Save/Get link")
    Rel_D(queue, worker, "Get raw files")
    Rel_D(worker, bucket_compress, "Save compressed files")
```

## :card_file_box: Database

```mermaid
  erDiagram
    folder {
        INT id
        INT parent_id
        STRING name
        DATETIME created_at
        DATETIME updated_at
        BOOL deleted
    }

    file {
        INT id
        INT folder_id
        INT owner_id
        STRING name
        STRING type
        STRING path
        DATETIME created_at
        DATETIME updated_at
        BOOL deleted
    }

    user {
        INT id
        STRING name
        STRING login
        STRING password
        DATETIME created_at
        DATETIME updated_at
        DATETIME last_login
        BOOL deleted
    }

    user ||--o{ file : "one to many"
    folder ||--o{ file : "one to many"
    folder ||--o{ folder : "one to many"
```

## :computer: Technologies

- [Go](https://golang.org/)
- [RabbitMQ](https://www.rabbitmq.com/)
- [Docker](https://www.docker.com/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## :scroll: Requirements

- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## :cd: Installation

```sh
git clone git@github.com:filipe1309/agl-go-driver.git
```

```sh
cd agl-go-driver
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## :runner: Running

```sh
make run
```

> Access http://localhost

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ## :white_check_mark: Tests

After up the container:

```sh
docker-compose exec -t {{ CONTAINER_SERVICE_NAME }} ./vendor/bin/phpunit
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate. -->

## :pin: Roadmap
- [ ] replace `aws-sdk-go` with `aws-sdk-go-v2`

license:
## :memo: License

[MIT](https://choosealicense.com/licenses/mit/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 🧙‍♂️ About Me

<p align="center">
    <a style="font-weight: bold" href="https://github.com/filipe1309/agl-go-driver/">
    <img style="border-radius:50%" width="100px; "src="https://github.com/filipe1309.png"/>
    </a>
</p>

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## :clap: Acknowledgments

- [ShubcoGen Template™](https://github.com/filipe1309/shubcogen-template)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

<p align="center">
    Done with&nbsp;&nbsp;♥️&nbsp;&nbsp;by <a style="font-weight: bold" href="https://github.com/filipe1309/">filipe1309</a> 🖖
</p>

---

