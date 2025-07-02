# [agree:user:pydantic]
class UserSchema(BaseModel):
    id: int
    username: str
    email: EmailStr
    full_name: str | None = None

    class Config:
        orm_mode = True


# [agree:end]
#


# [agree:post:pydantic]
class PostSchema(BaseModel):
    id: int
    user: str


# [agree:end]
