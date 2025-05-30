from pydantic import BaseModel, EmailStr

class UserSchema(BaseModel):
    id: int
    username: str
    email: EmailStr
    full_name: str | None = None

    class Config:
        orm_mode = True

class PostSchema(BaseModel):
    id: int
    title: str
    content: str | None = None
    author_id: int

class AddressSchema(BaseModel):
    id: int
    user: UserSchema
    street: str
    city: str
    state: str
    zip_code: str

class OrganizationSchema(BaseModel):
    id: int
    name: str
    domain: str
    description: str | None = None
    owner: UserSchema
    departments: list[str]

    class Config:
        orm_mode = True
