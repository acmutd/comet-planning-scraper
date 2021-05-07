import chevrotain from 'chevrotain';

const createToken = chevrotain.createToken;
const tokenMatcher = chevrotain.tokenMatcher;
const Lexer = chevrotain.Lexer;
const EmbeddedActionsParser = chevrotain.EmbeddedActionsParser;

// define the base cases, which are the components that make up expressions
const And = createToken({ name: 'And', pattern: /and/ });
const Or = createToken({ name: 'Or', pattern: /or/ });

const LParen = createToken({ name: 'LParen', pattern: /\(/ });
const RParen = createToken({ name: 'RParen', pattern: /\)/ });

const Course = createToken({ name: 'Course', pattern: /([A-Z]+ [0-9][A-Z0-9][0-9]+)/ });
const RandomRequest = createToken({ name: 'RandomRequest', pattern: /((?! and | or |\(|\)).)+/ });
const Grade = createToken({
  name: 'Grade',
  pattern: /with a (?:minimum )*grade (of )*[ABC]-*\+*( or (?:higher|better))*/,
});

const WhiteSpace = createToken({
  name: 'WhiteSpace',
  pattern: /\s+/,
  group: Lexer.SKIPPED,
});

const Comma = createToken({ name: 'Comma', pattern: /,/, group: Lexer.SKIPPED });

// define all tokens with their order of precedence
const allTokens = [
  WhiteSpace,
  Comma,
  And,
  Or,
  Course,
  Grade,
  RandomRequest,
  LParen,
  RParen, //, IrrelevantWord
];

// create lexer instance
const CalculatorLexer = new Lexer(allTokens);

// combine array of expressions into object with "and" key
// used to generate intuitive prerequisite graph
function generateAnd(children) {
  if (children.length > 1) return { and: children };
  else return children;
}

// combine array of expressions into object with "or" key
function generateOr(children) {
  if (children.length > 1) return { or: children };
  else return children;
}

class Calculator extends EmbeddedActionsParser {
  constructor() {
    super(allTokens);
    const $ = this;

    $.RULE('expression', () => {
      let res = $.SUBRULE($.andExpression);
      return res;
    });

    $.RULE('andExpression', () => {
      let value: any = []; // TODO
      // parsing part
      value.push($.SUBRULE($.orExpression));
      $.MANY(() => {
        // consuming 'AdditionOperator' will consume
        // either Plus or Minus as they are subclasses of AdditionOperator
        $.CONSUME(And);
        //  the index "2" in SUBRULE2 is needed to identify the unique
        // position in the grammar during runtime
        value.push($.SUBRULE2($.orExpression));
      });

      return value.length === 1 ? value[0] : generateAnd(value);
    });

    $.RULE('orExpression', () => {
      let value: any = [];

      // parsing part
      value.push($.SUBRULE($.atomicBooleanExpression));
      $.MANY(() => {
        $.CONSUME(Or);
        let val = $.SUBRULE2($.atomicBooleanExpression);
        value.push(val);
      });

      return value.length === 1 ? value[0] : generateOr(value);
    });

    $.RULE('atomicBooleanExpression', () =>
      $.OR([
        // parenthesisExpression has the highest precedence and thus it
        // appears in the "lowest" leaf in the expression ParseTree.
        { ALT: () => $.SUBRULE($.parenthesisExpression) },
        { ALT: () => $.SUBRULE($.courseExpression) },
        {
          ALT: () => {
            let rand = $.CONSUME(RandomRequest).image;
            return { course: rand, type: 'special' };
          },
        },
      ]),
    );

    $.RULE('courseExpression', () => {
      let course;
      let grade: any = -1;

      course = $.CONSUME(Course);
      $.OPTION(() => {
        grade = $.CONSUME(Grade);
      });

      if (grade == -1) return { course: course.image, grade: '' };
      return { course: course.image, grade: grade.image };
    });

    $.RULE('parenthesisExpression', () => {
      let expValue;
      let grade = '';

      $.CONSUME(LParen);
      expValue = $.SUBRULE($.andExpression);
      $.CONSUME(RParen);

      $.OPTION(() => {
        grade = $.CONSUME(Grade).image;
      });
      let res = { courses: expValue, grade: grade };
      return res;
    });

    // very important to call this after all the rules have been defined.
    // otherwise the parser may not work correctly as it will lack information
    // derived during the self analysis phase.
    this.performSelfAnalysis();
  }
}

const parser = new Calculator();

export function parseInput(text) {
  const lexingResult = CalculatorLexer.tokenize(text);
  // "input" is a setter which will reset the parser's state.
  parser.input = lexingResult.tokens;
  let res = parser.expression();

  if (parser.errors.length > 0) {
    throw new Error('sad sad panda, Parsing errors detected');
  }
  return res;
}

export function prettyPrint(text) {
  console.log(text);
  let res = parseInput(text);
  console.log(JSON.stringify(res, null, 2));
  return res;
}
