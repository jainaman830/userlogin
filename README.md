# userlogin
user login and registration
# Register api
Type : POST
URL : "/user/register"
sample input : {
    "username":"jainaman2",
    "firstname":"aman",
    "lastname":"jain",
    "email":"test1@gmal.com",
    "password":"Password.1"
}
sample output : 
{
    "message": "A verification mail has been sent to your registered mail."
}

# login api
Type : POST
URL : "/user/login"
Sample input :
{
    "Username":"jainaman",
    "Password":"password"
}
Sample output :
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDc5MjE2NDcsImlkIjoiNjVjY2NlMzU3NTc1NmM0ZDllMzBlMGI1IiwidXNlcm5hbWUiOiJqYWluYW1hbiJ9.qtLLzcj_0LDRcUfSZjEa5QPs2Z1D7WOT8Q3GMGr4e6s",
    "user": {
        "id": "65ccce3575756c4d9e30e0b5",
        "username": "jainaman",
        "firstname": "aman",
        "lastname": "jain",
        "email": "test@gmail.com",
        "password": "password",
        "CreatedOn": "2024-02-14T14:29:09.465Z"
    }
}

# userinfo api
Type : GET
URL :"/user/userinfo"
Header : "Authorization"
Output : {
    "id": "65ccce3575756c4d9e30e0b5",
    "username": "jainaman",
    "firstname": "",
    "lastname": "",
    "email": "<nil>",
    "password": "",
    "CreatedOn": "0001-01-01T00:00:00Z"
} 