import asyncio
import random
import string
import time

from fastapi import FastAPI, Depends, HTTPException, Header
from starlette.requests import Request
from fastapi.security import http
from starlette.responses import JSONResponse

import databases
import orm
import sqlalchemy

database = databases.Database("sqlite:///db.sqlite")
metadata = sqlalchemy.MetaData()
app = FastAPI()


class User(orm.Model):
    __tablename__ = "user"
    __metadata__ = metadata
    __database__ = database

    id = orm.Integer(primary_key=True)
    name = orm.String(max_length=100, unique=True, index=True)
    language = orm.String(max_length=10)
    is_active = orm.Boolean()

    @classmethod
    async def create(cls, name: str, **kwargs):
        return await cls.objects.create(name=name, **kwargs)


class UserRelations(orm.Model):
    __tablename__ = "relations"
    __metadata__ = metadata
    __database__ = database

    id = orm.Integer(primary_key=True)
    owner = orm.ForeignKey(User)
    target = orm.ForeignKey(User)
    kind = orm.String(max_length=10, index=True, default="friend")


def get_current_user_id(token=Depends(http.HTTPBearer())):
    return 1


async def get_current_user(user_id: str = Depends(get_current_user_id)):
    return await User.objects.get(id=user_id)


@app.middleware("http")
async def add_process_time_header(request: Request, call_next):
    start_time = time.perf_counter()
    response = await call_next(request)
    process_time = time.perf_counter() - start_time
    response.headers["X-Process-Time"] = str(process_time)
    return response


@app.get("/")
async def root():
    await asyncio.sleep(1)
    return "".join([random.choice(string.ascii_letters) for _ in range(10)])



@app.get("/me")
async def get_me(user: User = Depends(get_current_user)):
    friends = await database.fetch_all(
        "SELECT name, relations.kind FROM user JOIN relations ON relations.target = user.id WHERE relations.owner = :owner",
        {"owner": user.id},
    )
    return {"name": user.name, "friends": friends}


@app.get("/rank")
async def get_rang(accept_language: str = Header("en"), user_agent: str = Header("No")):
    rank = 0
    for user in await User.objects.filter(language=accept_language).all():
        for relation in await UserRelations.objects.filter(owner=user).all():
            rank += {"friend": 10, "dude": 3, "relative": 15}[relation.kind]
    return {"rank": rank, "language": accept_language}


async def create_data():
    user = await User.create("test", is_active=True, language="en")
    users = [
        user,
    ]
    for i in range(10000):
        name = "".join([random.choice(string.ascii_letters) for _ in range(10)])
        language = random.choice(["en", "es", "th"])
        users.append(await User.create(name, is_active=True, language=language))
    for i in range(5000):
        await UserRelations.objects.create(
            owner=random.choice(users), target=random.choice(users), kind=random.choice(["friend", "relative", "dude"])
        )


if __name__ == "__main__":
    engine = sqlalchemy.create_engine(str(database.url))
    metadata.create_all(engine)
    asyncio.run(create_data())
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
