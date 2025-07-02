# Test data with intentionally mismatched schemas for testing

# [agree:mismatch_test:pydantic]
class MismatchTestSchema(BaseModel):
    id: int
    name: str
    email: str
    age: int  # This field missing in TypeScript
    # missing 'score' field that exists in TypeScript
# [agree:end]

# [agree:mismatch_test:sqlalchemy]
class MismatchTest(Base):
    __tablename__ = "mismatch_tests"
    id = Column(Integer, primary_key=True)
    name = Column(String, nullable=False) 
    email = Column(String, unique=True)
    age = Column(Integer)  # This field missing in TypeScript
    # missing 'score' field that exists in TypeScript
# [agree:end]