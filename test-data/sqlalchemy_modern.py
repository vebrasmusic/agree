# Modern SQLAlchemy 2.0+ syntax with Mapped[] types
from typing import Optional
from sqlalchemy import String, Integer, DateTime
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, relationship


class Base(DeclarativeBase):
    pass


# [agree:post:sqlalchemy]
class Post(Base):
    __tablename__ = "posts"

    id: Mapped[int] = mapped_column(Integer, primary_key=True)
    title: Mapped[str] = mapped_column(String(200), nullable=False)
    content: Mapped[Optional[str]] = mapped_column(String(5000))
    author_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id"))

    # Relationship
    author: Mapped["User"] = relationship("User", back_populates="posts")


# [agree:end]

