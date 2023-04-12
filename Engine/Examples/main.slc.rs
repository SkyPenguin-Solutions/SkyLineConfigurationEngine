ENGINE(true) {
    INIT true {
     constant PARSER_ERROR_CODE_TO_MISSING_SEMICOLON = 12;
     constant PARSER_ERROR_CODE_TO_MISSING_BRACKET   = 12;
     constant PARSER_ERROR_CODE_TO_MISSING_CODELINE  = 12;
 
     set DebugValue = true;
     set VerbsBalue = true;
     set Depth      = 1;
     set Color      = false;
 
     system("errors")   -> [Color, Depth, DebugValue];
   }
 };