# Test schemas that should be perfectly equivalent cross-language

# [agree:perfect:pydantic]
class PerfectSchema(BaseModel):
    id: int
    name: str
    active: bool
    score: float
    email: EmailStr
# [agree:end]

# [agree:perfect:sqlalchemy]
class Perfect(Base):
    __tablename__ = "perfect"
    id = Column(Integer, primary_key=True)
    name = Column(String, nullable=False)
    active = Column(Boolean, default=True)
    score = Column(Float)
    email = Column(String)  # Will map to string, not email
# [agree:end]