"""Utility functions and mappings for parser"""

# SQLAlchemy type to Python type mapping
SQLALCHEMY_TYPE_MAP = {
    # Integer types
    "Integer": "int",
    "BigInteger": "int",
    "SmallInteger": "int",
    
    # String types
    "String": "str",
    "Text": "str",
    "Unicode": "str",
    "UnicodeText": "str",
    "VARCHAR": "str",
    "CHAR": "str",
    
    # Float types
    "Float": "float",
    "Numeric": "float",
    "DECIMAL": "float",
    "REAL": "float",
    
    # Boolean
    "Boolean": "bool",
    
    # DateTime types
    "DateTime": "datetime",
    "Date": "date",
    "Time": "time",
    "TIMESTAMP": "datetime",
    
    # Binary types
    "LargeBinary": "bytes",
    "BLOB": "bytes",
    
    # JSON types
    "JSON": "dict",
    "JSONB": "dict",
}


def map_sqlalchemy_type(sqlalchemy_type: str) -> str:
    """
    Maps a SQLAlchemy type to its Python equivalent.
    
    Args:
        sqlalchemy_type: The SQLAlchemy type name (e.g., 'Integer', 'String')
        
    Returns:
        The corresponding Python type name (e.g., 'int', 'str')
    """
    return SQLALCHEMY_TYPE_MAP.get(sqlalchemy_type, sqlalchemy_type)
