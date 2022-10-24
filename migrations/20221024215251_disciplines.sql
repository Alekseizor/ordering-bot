-- +goose Up
-- +goose StatementBegin
CREATE TABLE disciplines
(
    id      serial PRIMARY KEY ,
    name    varchar not null
);

INSERT INTO disciplines
VALUES (1,'MATLAB');
INSERT INTO disciplines (name)
VALUES ('MS Office(Word, Excel, Access)');
INSERT INTO disciplines (name)
VALUES ('Mathcad');
INSERT INTO disciplines (name)
VALUES ('Аналитическая геометрия');
INSERT INTO disciplines (name)
VALUES ('Английский язык');
INSERT INTO disciplines (name)
VALUES ('Архитектура');
INSERT INTO disciplines (name)
VALUES ('Бухгалтерский учет и анализ');
INSERT INTO disciplines (name)
VALUES ('Водоснабжение и водоотведение');
INSERT INTO disciplines (name)
VALUES ('Геодезия');
INSERT INTO disciplines (name)
VALUES ('Геология');
INSERT INTO disciplines (name)
VALUES ('Дискретная математика');
INSERT INTO disciplines (name)
VALUES ('Инженерная и компьютерная графика');
INSERT INTO disciplines (name)
VALUES ('Интегралы и дифференциальные уравнения');
INSERT INTO disciplines (name)
VALUES ('Информатика');
INSERT INTO disciplines (name)
VALUES ('Искусствоведение');
INSERT INTO disciplines (name)
VALUES ('История');
INSERT INTO disciplines (name)
VALUES ('Кратные интегралы и ряды');
INSERT INTO disciplines (name)
VALUES ('Культурология');
INSERT INTO disciplines (name)
VALUES ('Линейная алгебра');
INSERT INTO disciplines (name)
VALUES ('Математика');
INSERT INTO disciplines (name)
VALUES ('Математический анализ');
INSERT INTO disciplines (name)
VALUES ('Материаловедение');
INSERT INTO disciplines (name)
VALUES ('Метрология');
INSERT INTO disciplines (name)
VALUES ('Механика грунтов');
INSERT INTO disciplines (name)
VALUES ('Механика жидкости и газа (Гидравлика)');
INSERT INTO disciplines (name)
VALUES ('Начертательная геометрия');
INSERT INTO disciplines (name)
VALUES ('Политология');
INSERT INTO disciplines (name)
VALUES ('Правоведение');
INSERT INTO disciplines (name)
VALUES ('Психология');
INSERT INTO disciplines (name)
VALUES ('Сопротивление материалов');
INSERT INTO disciplines (name)
VALUES ('Социология');
INSERT INTO disciplines (name)
VALUES ('Строительная механика');
INSERT INTO disciplines (name)
VALUES ('Строительные материалы');
INSERT INTO disciplines (name)
VALUES ('Строительные машины');
INSERT INTO disciplines (name)
VALUES ('Теоретическая механика');
INSERT INTO disciplines (name)
VALUES ('Теория вероятностей');
INSERT INTO disciplines (name)
VALUES ('Теория механизмов и машин');
INSERT INTO disciplines (name)
VALUES ('Теория поля');
INSERT INTO disciplines (name)
VALUES ('Теория функции комплексных переменных');
INSERT INTO disciplines (name)
VALUES ('Теория функции нескольких переменных');
INSERT INTO disciplines (name)
VALUES ('Технология процессов в строительстве');
INSERT INTO disciplines (name)
VALUES ('Уравнения математической физики');
INSERT INTO disciplines (name)
VALUES ('Физика');
INSERT INTO disciplines (name)
VALUES ('Философия');
INSERT INTO disciplines (name)
VALUES ('Химия');
INSERT INTO disciplines (name)
VALUES ('Экономика');
INSERT INTO disciplines (name)
VALUES ('Электродинамика и распространение радиоволн');
INSERT INTO disciplines (name)
VALUES ('Электротехника');
INSERT INTO disciplines (name)
VALUES ('Детали машин');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
