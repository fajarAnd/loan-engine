## 🛠 Tech Stack

### **Backend Framework**
- **Language**: Go (Golang)
- **HTTP Framework**: Chi Router (RESTful API)
- **Design Pattern**: Dependency Injection

### **Database**
- **Primary DB**: PostgreSQL
- **ORM**: PGX (v5.7.5) - High-performance PostgreSQL driver
- **Migrations**: golang-migrate

---

### ERD
![erd.png](diagram/output/erd.png)

---

## API Spec

### API Endpoints
| No. | Fungsionalitas                  | HTTP Method | Endpoint                                 |
| :-- | :------------------------------ | :---------- | :--------------------------------------- |
| 1.  | Login User                      | `POST`      | `/api/v1/login`                          |
| 2.  | Create Loan Proposal            | `POST`      | `/api/v1/loans`                          |
| 3.  | Get Loan Details                | `GET`       | `/api/v1/loans/{id}`                     |
| 4.  | Approve Loan                    | `PUT`       | `/api/v1/loans/{id}/approve`             |
| 5.  | Disburse Loan                   | `PUT`       | `/api/v1/loans/{id}/disburse`            |
| 6.  | List Available Loans for Investment | `GET`       | `/api/v1/loans/available`                |
| 7.  | Make Investment in Loan         | `POST`      | `/api/v1/loans/{id}/investments`         |
| 8.  | Get Investor's Investment Portfolio | `GET`       | `/api/v1/investors/{investor_id}/portfolio` |
| 9.  | Upload Document Files           | `POST`      | `/api/v1/files/upload`                   |
| 10. | Download/View Document          | `GET`       | `/api/v1/files/{file_id}`                |
| 11. | Basic Health Check              | `GET`       | `/api/v1/__health`                       |
