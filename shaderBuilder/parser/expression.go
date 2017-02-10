package parser

import (
	"math"
	"strconv"
)

type Variable struct {
	Name  string
	Value float64
}

type Expression interface {
	Evaluate(variables ...Variable) float64
}

type BinaryExpression struct {
	Operator Token
	Left     Expression
	Right    Expression
}

func (b BinaryExpression) Evaluate(variables ...Variable) float64 {
	switch b.Operator {
	case PLUS:
		return b.Left.Evaluate(variables...) + b.Right.Evaluate(variables...)
	case MINUS:
		return b.Left.Evaluate(variables...) - b.Right.Evaluate(variables...)
	case MULTIPLY:
		return b.Left.Evaluate(variables...) * b.Right.Evaluate(variables...)
	case SLASH:
		return b.Left.Evaluate(variables...) / b.Right.Evaluate(variables...)
	case REMAINDER:
		leftInt, rightInt := int(b.Left.Evaluate(variables...)), int(b.Right.Evaluate(variables...))
		return float64(leftInt % rightInt)
	case POWER:
		return math.Pow(b.Left.Evaluate(variables...), b.Right.Evaluate(variables...))
	default:
	}
	return 0
}

type Identifier struct {
	Name string
}

func (i Identifier) Evaluate(variables ...Variable) float64 {
	for _, variable := range variables {
		if variable.Name == i.Name {
			return variable.Value
		}
	}
	return 0
}

type Literal struct {
	Value float64
}

func (l Literal) Evaluate(variables ...Variable) float64 {
	return l.Value
}

func (p *Parser) parseSignedExpression() Expression {
	var left Expression = Literal{Value: 0}

	if p.token == PLUS || p.token == MINUS {
		if p.token == MINUS {
			left = Literal{-1}
		} else {
			left = Literal{1}
		}
		p.nextNoWhitespace()

		left = BinaryExpression{
			Operator: MULTIPLY,
			Left:     left,
			Right:    p.parseExpression(),
		}
	}

	return left
}

func (p *Parser) parseGroupExpression() Expression {
	exp := p.parseSignedExpression()

	if p.token == LEFT_PARENTHESIS {
		p.nextNoWhitespace()
		exp = p.parseExpression()

		if p.token != RIGHT_PARENTHESIS {
			p.unexpectedError()
		}
		p.nextNoWhitespace()
	}

	return exp
}

func (p *Parser) parseIdentifierExpression() Expression {
	exp := p.parseGroupExpression()

	if p.token == IDENTIFIER {
		exp = Identifier{Name: p.literal}
		p.nextNoWhitespace()
	}

	return exp
}

func (p *Parser) parseLiteralExpression() Expression {
	exp := p.parseIdentifierExpression()

	if p.token == NUMBER {
		value, _ := strconv.ParseFloat(p.literal, 64)
		exp = Literal{Value: value}
		p.nextNoWhitespace()
	}

	return exp
}

func (p *Parser) parsePowerExpression() Expression {
	nextParse := p.parseLiteralExpression
	left := nextParse()

	for p.token == POWER {
		tkn := p.token
		p.nextNoWhitespace()

		left = BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    nextParse(),
		}
	}

	return left
}

func (p *Parser) parseMultiplicativeExpression() Expression {
	nextParse := p.parsePowerExpression
	left := nextParse()

	for p.token == MULTIPLY || p.token == SLASH || p.token == REMAINDER {
		tkn := p.token
		p.nextNoWhitespace()

		left = BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    nextParse(),
		}
	}

	return left
}

func (p *Parser) parseAdditiveExpression() Expression {
	nextParse := p.parseMultiplicativeExpression
	left := nextParse()

	for p.token == PLUS || p.token == MINUS {
		tkn := p.token
		p.nextNoWhitespace()

		left = BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    nextParse(),
		}
	}

	return left
}

func (p *Parser) parseExpression() Expression {
	return p.parseAdditiveExpression()
}
