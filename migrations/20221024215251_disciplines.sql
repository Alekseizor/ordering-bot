-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS disciplines
(
    id      serial PRIMARY KEY ,
    name    varchar not null
);

INSERT INTO disciplines (name)
VALUES ('MATLAB');
INSERT INTO disciplines (name)
VALUES ('MS Office(Word, Excel, Access)');
INSERT INTO disciplines (name)
VALUES ('Mathcad');
INSERT INTO disciplines (name)
VALUES ('Аналитическая геометрия');
INSERT INTO disciplines (name)
VALUES ('Английский язык');
INSERT INTO disciplines (name)
VALUES ('Детали машин');
INSERT INTO disciplines (name)
VALUES ('Дискретная математика');
INSERT INTO disciplines (name)
VALUES ('Инженерная и компьютерная графика');
INSERT INTO disciplines (name)
VALUES ('Интегралы и дифференциальные уравнения');
INSERT INTO disciplines (name)
VALUES ('Информатика');
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
VALUES ('Менеджмент');
INSERT INTO disciplines (name)
VALUES ('Метрология');
INSERT INTO disciplines (name)
VALUES ('Механика жидкости и газа');
INSERT INTO disciplines (name)
VALUES ('Начертательная геометрия');
INSERT INTO disciplines (name)
VALUES ('Организация производства');
INSERT INTO disciplines (name)
VALUES ('Основы конструирования приборов');
INSERT INTO disciplines (name)
VALUES ('Основы теории цепей');
INSERT INTO disciplines (name)
VALUES ('Основы технологии приборостроения');
INSERT INTO disciplines (name)
VALUES ('Политология');
INSERT INTO disciplines (name)
VALUES ('Правоведение');
INSERT INTO disciplines (name)
VALUES ('Практика');
INSERT INTO disciplines (name)
VALUES ('Прикладная статистика');
INSERT INTO disciplines (name)
VALUES ('Психология');
INSERT INTO disciplines (name)
VALUES ('Системный анализ и принятие решений');
INSERT INTO disciplines (name)
VALUES ('Сопротивление материалов');
INSERT INTO disciplines (name)
VALUES ('Социология');
INSERT INTO disciplines (name)
VALUES ('Теоретическая механика');
INSERT INTO disciplines (name)
VALUES ('Теоретические основы электротехники');
INSERT INTO disciplines (name)
VALUES ('Теория вероятностей');
INSERT INTO disciplines (name)
VALUES ('Теория механизмов и машин');
INSERT INTO disciplines (name)
VALUES ('Теория поля');
INSERT INTO disciplines (name)
VALUES ('Теория функции комплексных переменных и операционное исчисление');
INSERT INTO disciplines (name)
VALUES ('Теория функции нескольких переменных');
INSERT INTO disciplines (name)
VALUES ('Термодинамика');
INSERT INTO disciplines (name)
VALUES ('Технология конструкционных материалов');
INSERT INTO disciplines (name)
VALUES ('Уравнения математической физики');
INSERT INTO disciplines (name)
VALUES ('Физика');
INSERT INTO disciplines (name)
VALUES ('Физкультура');
INSERT INTO disciplines (name)
VALUES ('Философия');
INSERT INTO disciplines (name)
VALUES ('Финансирование инновационной деятельности');
INSERT INTO disciplines (name)
VALUES ('Химия');
INSERT INTO disciplines (name)
VALUES ('Цифровые устройства и микропроцессоры');
INSERT INTO disciplines (name)
VALUES ('Экономика');
INSERT INTO disciplines (name)
VALUES ('Электроника');
INSERT INTO disciplines (name)
VALUES ('Электротехника');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS disciplines;
-- +goose StatementEnd
