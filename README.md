<p>A package of validators and sanitizers for strings, structs and collections.</p>
<p>features:</p>
<ul>
  <li>Customizable Attributes.</li>
  <li>Customizable error messages.</li>
  <li>Support i18n messages</li>
</ul>
<h2>Installation</h2>
<p>Make sure that Go is installed on your computer. Type the following command in your terminal:</p>
<p>go get github.com/syssam/go-validator</p>
<h2>Usage and documentation</h2>
<h5>Examples:</h5>
<ul>
  <li><a href="https://github.com/syssam/go-validator/tree/master/examples/simple">Simple</a></li>
  <li><a href="https://github.com/syssam/go-validator/tree/master/examples/translations">Translations</a></li>
  <li><a href="https://github.com/syssam/go-validator/tree/master/examples/gin">Gin</a></li>
  <li><a href="https://github.com/syssam/go-validator/tree/master/examples/custom">Custom Validation Rules</a></li>
</ul>
<h2>Available Validation Rules</h2>
<ul>
    <li><a>required</a></li>
    <li><a>requiredIf</a></li>
    <li><a>requiredUnless</a></li>
    <li><a>requiredWith</a></li>
    <li><a>requiredWithAll</a></li>
    <li><a>requiredWithout</a></li>
    <li><a>requiredWithoutAll</a></li>
    <li><a>email</a></li>
    <li><a>between</a></li>
    <li><a>size</a></li>
    <li><a>max</a></li>
    <li><a>min</a></li>
    <li><a>gt</a></li>
    <li><a>gte</a></li>
    <li><a>lt</a></li>
    <li><a>lte</a></li>
</ul>
<h4 id="rule-required">required</h4>
<p>The field under validation must be present in the input data and not empty. A field is considered "empty" if one of the following conditions are true:</p>
<div class="content-list">
  <ul>
    <li>The value is <code class="language-php"><span class="token keyword">nil</span></code>.</li>
    <li>The value is an empty string.</li>
    <li>The value is an empty array | map</li>
  </ul>
</div>
<h4 id="rule-requiredIf">requiredIf=anotherfield|value|...</h4>
<p>The field under validation must be present and not empty if the anotherfield field is equal to any value.</p>
<h4 id="rule-requiredIf">requiredUnless=anotherfield|value|...</h4>
<p>The field under validation must be present and not empty unless the anotherfield field is equal to any value.</p>
<h4 id="rule-requiredIf">requiredWith=anotherfield|anotherfield|...</h4>
<p>The field under validation must be present and not empty only if any of the other specified fields are present.</p>
<h4 id="rule-requiredIf">requiredWithAll=anotherfield|anotherfield|...</h4>
<p>The field under validation must be present and not empty only if all of the other specified fields are present.</p>
<h4 id="rule-requiredIf">requiredWithout=anotherfield|anotherfield|...</h4>
<p>The field under validation must be present and not empty only when any of the other specified fields are not present.</p>
<h4 id="rule-requiredIf">requiredWithoutAll=anotherfield|anotherfield|...</h4>
<p>The field under validation must be present and not empty only when all of the other specified fields are not present.</p>
<h4 id="rule-email">email</h4>
<p>The field under validation must be formatted as an e-mail address.</p>
<h4 id="rule-between">between=min|max</h4>
<p>The field under validation must have a size between the given min and max. String, Number, Array, Map are evaluated in the same fashion as the size rule.</p>
<h4 id="rule-max">size=value</h4>
<p>The field under validation must have a size matching the given value. For string data, value corresponds to the number of characters. For numeric data, value corresponds to a given integer value. For an array | map | slice, size corresponds to the count of the array | map | slice.</p>
<h4 id="rule-max">max=value</h4>
<p>The field under validation must be less than or equal to a maximum value. String, Number, Array, Map are evaluated in the same fashion as the size rule.</p>
<h4 id="rule-min">min=value</h4>
<p>The field under validation must be greater than or equal to a minimum value. String, Number, Array, Map are evaluated in the same fashion as the size rule.</p>
<h4 id="rule-gt">gt=anotherfield</h4>
<p>The field under validation must be greater than the given field. The two fields must be of the same type. String, Number, Array, Map are evaluated using the same conventions as the size rule.</p>
<h4 id="rule-gte">gte=anotherfield</h4>
<p>The field under validation must be greater than or equal to the given field. The two fields must be of the same type. String, Number, Array, Map are evaluated using the same conventions as the size rule.</p>
<h4 id="rule-lt">lt=anotherfield</h4>
<p>The field under validation must be less than the given field. The two fields must be of the same type. String, Number, Array, Map are evaluated using the same conventions as the size rule.</p>
<h4 id="rule-lte">lte=anotherfield</h4>
<p>The field under validation must be less than or equal to the given field. The two fields must be of the same type. String, Number, Array, Map are evaluated using the same conventions as the size rule.</p>
<h2>Custom Validation Rules</h2>
<div class="highlight highlight-source-go">
  <pre>
  validator.CustomTypeTagMap.Set("customValidator", func CustomValidator(v reflect.Value, o reflect.Value, validTag *validator.ValidTag) bool {
    return false
  })
  </pre>
</div>
<h2>List of functions:</h2>
<div class="highlight highlight-source-go">
  <pre>
    IsRequiredIf(v reflect.Value, anotherfield reflect.Value, params ...string)
    IsIn(str string, params ...string) bool 
    IsEmail(str string) bool
    Between(v reflect.Value, params ...string) bool
  </pre>
</div>