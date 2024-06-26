Этот проект состоит из серии упражнений по работе с бинарными деревьями и кучами. Ниже приведен краткий обзор каждого упражнения вместе с их реализацией.

## Упражнение 00: Игрушки на дереве


**Цель**: Создать функцию `areToysBalanced`, чтобы определить, равное ли количество игрушек (значений true) в левой и правой поддеревьях заданного бинарного дерева.

**Описание**: Функция принимает корень дерева и возвращает true, если количество игрушек в левой и правой части одинаково, и false в противном случае.

## Упражнение 01: Украшение

**Название алгоритма**: Зигзагообразный обход в ширину бинарного дерева.

**Цель**: Написать функцию `unrollGarland`, которая обходит бинарное дерево слой за слоем в зигзагообразном порядке и возвращает массив булевых значений, представляющих узлы дерева.

**Описание**: Функция выполняет обход дерева по уровням, меняя направление обхода на каждом уровне (справа налево и слева направо), и возвращает массив значений узлов.

## Упражнение 02: Куча подарков

Реализация структуры данных "Куча" для управления подарками.

**Цель**: Реализовать структуру данных кучи для управления подарками, приоритет которых определяется их ценностью и размером. Создать функцию `getNCoolestPresents`, которая возвращает N самых крутых подарков.

Функция использует кучу для хранения и извлечения самых ценных подарков, сортируя их по убыванию ценности и размера.

## Упражнение 03: Рюкзак

**Название алгоритма**: Решение задачи о рюкзаке.

**Цель**: Реализовать алгоритм решения задачи о рюкзаке, чтобы выбрать самые ценные подарки, которые поместятся в заданный объем.

**Описание**: Функция принимает массив подарков и вместимость рюкзака, возвращая массив подарков с максимальной суммарной ценностью, которые можно уместить в рюкзаке.