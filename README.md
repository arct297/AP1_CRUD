# Clinic Management System 
Purpose: to book doctor appointments, manage medical records, and schedule.
## Target Audience:
- **Doctors**: To manage appointments, view patient records, and track treatment schedules.
- **Patients**: To book appointments, view medical records, and access their appointment history.

## Team Members:
- **Arseniy Sinyov**
- **Kumissay Zhalmagambetova**
<h2>Instructions:</h2>

#### 1.1 Install Go Lang:
Follow the installation guide on the [official Go website](https://golang.org/doc/install).

#### 1.2 Install PostgreSQL:
1. Download PostgreSQL from the [official website](https://www.postgresql.org/download/).
2. Follow the installation instructions for your operating system.

#### 1.3 Install Dependencies:
In the project directory, run the following commands to install the required Go modules:
```bash
go get github.com/gorilla/mux
go get gorm.io/driver/postgres
go get gorm.io/gorm
```

#### 1.4 Setting Up the Database:
After installing PostgreSQL, you need to set up the database and configure the connection.

1. **Create the Database**:
   Open your PostgreSQL terminal and create a new database for the project:
   ```sql
   CREATE DATABASE clinic_db;
   
   ```
   
2. **Create a User and a Table: You will need to create a user and assign privileges to the clinic_db database**:
After installing PostgreSQL, you need to set up the database and configure the connection.


   ```sql


    CREATE USER username WITH PASSWORD 'password';
    
    GRANT ALL PRIVILEGES ON DATABASE clinic_db TO username;

    CREATE TABLE patients (
        id SERIAL PRIMARY KEY,               
        name VARCHAR(100) NOT NULL,         
        age INT NOT NULL,                    
        gender VARCHAR(10) NOT NULL,         
        contact VARCHAR(50) NOT NULL,        
        address VARCHAR(255) NOT NULL,      
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
    );
  
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO username;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO username;
    GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO username;
    
    INSERT INTO patients (name, age, gender, contact, address) 
    VALUES ('Test Patient', 25, 'Female', '123-456-7890', '123 Main St');


#### 1.5 **Clone the repository**
If you haven't already, clone the repository to your local machine:
   ```bash
   git clone https://github.com/arct297/ClinicMS
   cd task2
  ```

#### 1.6 **Run the program**
   ```bash
   go run main.go
```

Tools: Go Lang, PostgreSQL, Postman and Python testing

