import psycopg2

try:
    conn = psycopg2.connect(
        dbname="magistrala",
        user="magistrala",
        password="magistrala",
        host="127.0.0.1",
        port="5432"
    )
    print("连接成功！")
except Exception as e:
    print("连接失败：", e)
