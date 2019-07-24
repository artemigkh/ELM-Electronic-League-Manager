from faker import Faker
from faker.providers import internet
from faker.providers import misc

fake = Faker()
fake.add_provider(internet)
fake.add_provider(misc)


class User:
    def __init__(self, t):
        self.email = fake.email()
        self.password = fake.password(length=10,
                                      special_chars=True,
                                      digits=True,
                                      upper_case=True,
                                      lower_case=True)
        r = t.http.post("http://localhost:8080/api/v1/users", json={
            "email": self.email,
            "password": self.password
        })

        if r.status_code != 201:
            print(r.json())
        t.assertEqual(201, r.status_code)
