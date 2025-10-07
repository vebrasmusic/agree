from datetime import datetime, timedelta
from typing import Optional, List, Dict, Union
from pydantic import BaseModel, Field
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column


# # Some arbitrary constants
# DEFAULT_LIMIT = 100
# API_VERSION = "v1.2.3"
# CONFIG: Dict[str, str] = {"env": "dev", "region": "us-west-1"}
#
#
# # Filler enum-like class
# class Status:
#     ACTIVE = "active"
#     INACTIVE = "inactive"
#     PENDING = "pending"
#
#
# # Filler utility function
# def calculate_offset(date: datetime, days: int = 7) -> datetime:
#     return date + timedelta(days=days)
#
#
# # Placeholder base class
# class Base(DeclarativeBase):
#     pass
#
#
# # Another filler model with no decorator
# class User(BaseModel):
#     id: int
#     username: str
#     email: Optional[str] = None


# AGREED schema (decorated)
@agree(target="Event")
class EventSchema(BaseModel):
    id: int
    date: datetime
    num_ppl: Optional[int] = Field(default=None, description="Number of people")
    test_col_1: Union[int, None] = None
    test_col_2: int | None = None


# # Another filler Pydantic model
# class LocationSchema(BaseModel):
#     id: int
#     address: str
#     city: str
#     country: str = "USA"
#
#
# # AGREED SQLAlchemy model (decorated)
@agree(target="Event")
class EventModel(Base):
    __tablename__ = "event"

    id: Mapped[int] = mapped_column(primary_key=True)
    date: Mapped[datetime]
    num_ppl: Mapped[Optional[int]]


@agree(target="Event")
class EventModel2(Base):
    __tablename__ = "event"

    id = Column(Integer, primary_key=True)
    date = Column(DateTime)
    num_ppl = Column(Integer, nullable=True)


#
#
# # Another filler ORM model
# @agree(target="location", fidelity=2)
# class LocationModel(Base):
#     __tablename__ = "location"
#
#     id: Mapped[int] = mapped_column(primary_key=True)
#     address: Mapped[str]
#     city: Mapped[str]
#     country: Mapped[str]
