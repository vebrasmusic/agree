# Test data with intentionally matching schemas for testing

# [agree:match_test:pydantic]
class MatchTestSchema(BaseModel):
    id: int
    name: str
    email: str
    active: bool
# [agree:end]

# [agree:match_test:sqlalchemy]
class MatchTest(Base):
    __tablename__ = "match_tests"
    id = Column(Integer, primary_key=True)
    name = Column(String, nullable=False)
    email = Column(String, unique=True)
    active = Column(Boolean, default=True)
# [agree:end]