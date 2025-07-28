import requests
import random
import string

URL = "http://localhost:8500/api/"

def gen_text(length: int) -> str:
    return ''.join(''.join(random.choices(string.ascii_uppercase, k=length)))

ROLE_GUEST = "guest"
ROLE_HOST = "host"
ROLE_ADMIN = "admin"

def register(username: str, password: str, role: str):
    return requests.post(
        URL + "register",
        json = {
            "username": username,
            "password": password,
            "email": f"{username}@gmail.com",
            "name": gen_text(6),
            "surname": gen_text(6),
            "role": role,
            "address": gen_text(10),
        }
    )

def login(username_or_email: str, password: str):
    return requests.post(
        URL + "login",
        json = {
            "usernameOrEmail": username_or_email,
            "password": password
        }
    )

def H(resp: requests.Response):
    if not resp.ok:
        print("Error", resp.status_code, resp.content)
    else:
        print(resp.json())


H(register("user1", "1234", ROLE_GUEST))
H(register("user1", "1234", ROLE_GUEST))

H(register("user2", "12345", ROLE_GUEST))

H(login("user1", "1234"))
H(login("user2", "1234"))
H(login("user1", "12345"))
H(login("user2", "12345"))